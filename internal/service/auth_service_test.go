package service

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"errors"
	"fmt"
	"io"
	"reflect"
	"sort"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/gitsang/order/internal/config"
	"github.com/gitsang/order/internal/model"
	"github.com/gitsang/order/internal/repository"
)

type serviceTestStore struct {
	mu         sync.Mutex
	users      map[uuid.UUID]model.User
	categories map[uuid.UUID]model.Category
	products   map[uuid.UUID]model.Product
	orders     map[uuid.UUID]model.Order
	items      map[uuid.UUID]model.OrderItem
	errors     map[string]error
}

func newServiceTestStore() *serviceTestStore {
	return &serviceTestStore{
		users:      make(map[uuid.UUID]model.User),
		categories: make(map[uuid.UUID]model.Category),
		products:   make(map[uuid.UUID]model.Product),
		orders:     make(map[uuid.UUID]model.Order),
		items:      make(map[uuid.UUID]model.OrderItem),
		errors:     make(map[string]error),
	}
}

func newTestRepositories(t *testing.T, store *serviceTestStore) (*repository.UserRepository, *repository.ProductRepository, *repository.OrderRepository) {
	t.Helper()

	db := sql.OpenDB(serviceTestConnector{store: store})
	t.Cleanup(func() { _ = db.Close() })

	gormDB, err := gorm.Open(postgres.New(postgres.Config{
		Conn:                 db,
		PreferSimpleProtocol: true,
		WithoutReturning:     true,
	}), &gorm.Config{SkipDefaultTransaction: true})
	if err != nil {
		t.Fatalf("open test gorm db: %v", err)
	}

	return repository.NewUserRepository(gormDB), repository.NewProductRepository(gormDB), repository.NewOrderRepository(gormDB)
}

type serviceTestConnector struct{ store *serviceTestStore }

func (c serviceTestConnector) Connect(context.Context) (driver.Conn, error) {
	return &serviceTestConn{store: c.store}, nil
}

func (c serviceTestConnector) Driver() driver.Driver { return serviceTestDriver{} }

type serviceTestDriver struct{}

func (serviceTestDriver) Open(string) (driver.Conn, error) {
	return &serviceTestConn{store: newServiceTestStore()}, nil
}

type serviceTestConn struct{ store *serviceTestStore }

func (c *serviceTestConn) Prepare(query string) (driver.Stmt, error) {
	return &serviceTestStmt{conn: c, query: query}, nil
}
func (c *serviceTestConn) Close() error              { return nil }
func (c *serviceTestConn) Begin() (driver.Tx, error) { return serviceTestTx{}, nil }
func (c *serviceTestConn) BeginTx(context.Context, driver.TxOptions) (driver.Tx, error) {
	return serviceTestTx{}, nil
}
func (c *serviceTestConn) Ping(context.Context) error { return nil }
func (c *serviceTestConn) ExecContext(_ context.Context, query string, args []driver.NamedValue) (driver.Result, error) {
	return c.store.exec(query, namedValues(args))
}
func (c *serviceTestConn) QueryContext(_ context.Context, query string, args []driver.NamedValue) (driver.Rows, error) {
	return c.store.query(query, namedValues(args))
}

type serviceTestStmt struct {
	conn  *serviceTestConn
	query string
}

func (s *serviceTestStmt) Close() error  { return nil }
func (s *serviceTestStmt) NumInput() int { return -1 }
func (s *serviceTestStmt) Exec(args []driver.Value) (driver.Result, error) {
	return s.conn.store.exec(s.query, args)
}
func (s *serviceTestStmt) Query(args []driver.Value) (driver.Rows, error) {
	return s.conn.store.query(s.query, args)
}
func (s *serviceTestStmt) ExecContext(_ context.Context, queryArgs []driver.NamedValue) (driver.Result, error) {
	return s.conn.store.exec(s.query, namedValues(queryArgs))
}
func (s *serviceTestStmt) QueryContext(_ context.Context, queryArgs []driver.NamedValue) (driver.Rows, error) {
	return s.conn.store.query(s.query, namedValues(queryArgs))
}

type serviceTestTx struct{}

func (serviceTestTx) Commit() error   { return nil }
func (serviceTestTx) Rollback() error { return nil }

type serviceTestResult int64

func (r serviceTestResult) LastInsertId() (int64, error) { return 0, nil }
func (r serviceTestResult) RowsAffected() (int64, error) { return int64(r), nil }

type serviceTestRows struct {
	columns []string
	values  [][]driver.Value
	pos     int
}

func (r *serviceTestRows) Columns() []string { return r.columns }
func (r *serviceTestRows) Close() error      { return nil }
func (r *serviceTestRows) Next(dest []driver.Value) error {
	if r.pos >= len(r.values) {
		return io.EOF
	}
	copy(dest, r.values[r.pos])
	r.pos++
	return nil
}

func namedValues(args []driver.NamedValue) []driver.Value {
	values := make([]driver.Value, len(args))
	for i, arg := range args {
		values[i] = arg.Value
	}
	return values
}

func (s *serviceTestStore) exec(query string, args []driver.Value) (driver.Result, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	q := normalizeSQL(query)
	switch {
	case strings.HasPrefix(q, "insert into users"):
		if err := s.errors["user.create"]; err != nil {
			return nil, err
		}
		user := userFromColumns(extractInsertColumns(query), args)
		s.users[user.ID] = user
		return serviceTestResult(1), nil
	case strings.HasPrefix(q, "insert into products"):
		if err := s.errors["product.create"]; err != nil {
			return nil, err
		}
		product := productFromColumns(extractInsertColumns(query), args)
		s.products[product.ID] = product
		return serviceTestResult(1), nil
	case strings.HasPrefix(q, "insert into orders"):
		if err := s.errors["order.create"]; err != nil {
			return nil, err
		}
		order := orderFromColumns(extractInsertColumns(query), args)
		s.orders[order.ID] = order
		return serviceTestResult(1), nil
	case strings.HasPrefix(q, "insert into order_items"):
		item := orderItemFromColumns(extractInsertColumns(query), args)
		s.items[item.ID] = item
		return serviceTestResult(1), nil
	case strings.HasPrefix(q, "update products"):
		if err := s.errors["product.update"]; err != nil {
			return nil, err
		}
		return s.updateProduct(query, args)
	case strings.HasPrefix(q, "update orders"):
		if err := s.errors["order.update_status"]; err != nil {
			return nil, err
		}
		return s.updateOrderStatus(args)
	default:
		return nil, fmt.Errorf("unexpected exec query %q with args %v", query, args)
	}
}

func (s *serviceTestStore) query(query string, args []driver.Value) (driver.Rows, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	q := normalizeSQL(query)
	switch {
	case strings.HasPrefix(q, "select * from users"):
		return s.queryUsers(q, args), nil
	case strings.HasPrefix(q, "select * from products"):
		return s.queryProducts(q, args), nil
	case strings.HasPrefix(q, "select * from categories"):
		return s.queryCategories(args), nil
	case strings.HasPrefix(q, "select * from orders"):
		return s.queryOrders(q, args), nil
	case strings.HasPrefix(q, "select * from order_items"):
		return s.queryOrderItems(q, args), nil
	default:
		return nil, fmt.Errorf("unexpected query %q with args %v", query, args)
	}
}

func (s *serviceTestStore) updateProduct(query string, args []driver.Value) (driver.Result, error) {
	q := normalizeSQL(query)
	if strings.Contains(q, "set deleted_at") {
		id := mustUUID(args[len(args)-1])
		product, ok := s.products[id]
		if !ok {
			return serviceTestResult(0), nil
		}
		delete(s.products, id)
		product.DeletedAt.Valid = true
		return serviceTestResult(1), nil
	}

	columns := extractSetColumns(query)
	values := mapColumns(columns, args[:len(columns)])
	id := mustUUID(args[len(args)-1])
	product, ok := s.products[id]
	if !ok {
		return serviceTestResult(0), nil
	}
	applyProductValues(&product, values)
	s.products[id] = product
	return serviceTestResult(1), nil
}

func (s *serviceTestStore) updateOrderStatus(args []driver.Value) (driver.Result, error) {
	if len(args) < 2 {
		return nil, errors.New("missing order status update args")
	}
	var status string
	var id uuid.UUID
	for _, arg := range args {
		if value, ok := arg.(string); ok {
			if _, err := uuid.Parse(value); err == nil {
				id = mustUUID(value)
				continue
			}
			status = value
		}
	}
	if id == uuid.Nil || status == "" {
		return nil, fmt.Errorf("missing order status update values: %v", args)
	}
	order, ok := s.orders[id]
	if !ok {
		return serviceTestResult(0), nil
	}
	order.Status = status
	s.orders[id] = order
	return serviceTestResult(1), nil
}

func (s *serviceTestStore) queryUsers(q string, args []driver.Value) driver.Rows {
	columns := userColumns()
	for _, user := range s.users {
		if matchesUser(q, user, args) {
			return &serviceTestRows{columns: columns, values: [][]driver.Value{userValues(user)}}
		}
	}
	return &serviceTestRows{columns: columns}
}

func (s *serviceTestStore) queryProducts(q string, args []driver.Value) driver.Rows {
	columns := productColumns()
	var rows [][]driver.Value
	for _, product := range sortedProducts(s.products) {
		if matchesProduct(q, product, args) {
			rows = append(rows, productValues(product))
		}
	}
	return &serviceTestRows{columns: columns, values: rows}
}

func (s *serviceTestStore) queryCategories(args []driver.Value) driver.Rows {
	columns := categoryColumns()
	var rows [][]driver.Value
	for _, arg := range args {
		if category, ok := s.categories[mustUUID(arg)]; ok {
			rows = append(rows, categoryValues(category))
		}
	}
	return &serviceTestRows{columns: columns, values: rows}
}

func (s *serviceTestStore) queryOrders(q string, args []driver.Value) driver.Rows {
	columns := orderColumns()
	var rows [][]driver.Value
	for _, order := range sortedOrders(s.orders) {
		if matchesOrder(q, order, args) {
			rows = append(rows, orderValues(order))
		}
	}
	return &serviceTestRows{columns: columns, values: rows}
}

func (s *serviceTestStore) queryOrderItems(q string, args []driver.Value) driver.Rows {
	columns := orderItemColumns()
	var rows [][]driver.Value
	for _, item := range sortedOrderItems(s.items) {
		if matchesOrderItem(q, item, args) {
			rows = append(rows, orderItemValues(item))
		}
	}
	return &serviceTestRows{columns: columns, values: rows}
}

func normalizeSQL(query string) string {
	query = strings.ToLower(query)
	query = strings.ReplaceAll(query, "\"", "")
	return strings.Join(strings.Fields(query), " ")
}

func extractInsertColumns(query string) []string {
	start := strings.Index(query, "(")
	end := strings.Index(query[start:], ")")
	return splitColumnList(query[start+1 : start+end])
}

func extractSetColumns(query string) []string {
	lower := strings.ToLower(query)
	setStart := strings.Index(lower, " set ")
	whereStart := strings.Index(lower, " where ")
	if setStart == -1 || whereStart == -1 {
		return nil
	}
	assignments := strings.Split(query[setStart+5:whereStart], ",")
	columns := make([]string, 0, len(assignments))
	for _, assignment := range assignments {
		columns = append(columns, cleanColumn(strings.Split(assignment, "=")[0]))
	}
	return columns
}

func splitColumnList(raw string) []string {
	parts := strings.Split(raw, ",")
	columns := make([]string, 0, len(parts))
	for _, part := range parts {
		columns = append(columns, cleanColumn(part))
	}
	return columns
}

func cleanColumn(column string) string {
	column = strings.TrimSpace(column)
	column = strings.Trim(column, "\"")
	if dot := strings.LastIndex(column, "."); dot >= 0 {
		column = column[dot+1:]
	}
	return strings.Trim(column, "\"")
}

func mapColumns(columns []string, args []driver.Value) map[string]driver.Value {
	values := make(map[string]driver.Value, len(columns))
	for i, column := range columns {
		if i < len(args) {
			values[column] = args[i]
		}
	}
	return values
}

func userFromColumns(columns []string, args []driver.Value) model.User {
	values := mapColumns(columns, args)
	user := model.User{ID: uuid.New()}
	if v, ok := values["id"]; ok {
		user.ID = mustUUID(v)
	}
	user.Username, _ = values["username"].(string)
	user.Password, _ = values["password"].(string)
	user.Name, _ = values["name"].(string)
	user.Phone, _ = values["phone"].(string)
	user.Role, _ = values["role"].(string)
	return user
}

func productFromColumns(columns []string, args []driver.Value) model.Product {
	values := mapColumns(columns, args)
	product := model.Product{ID: uuid.New()}
	applyProductValues(&product, values)
	return product
}

func applyProductValues(product *model.Product, values map[string]driver.Value) {
	if v, ok := values["id"]; ok {
		product.ID = mustUUID(v)
	}
	if v, ok := values["category_id"]; ok {
		product.CategoryID = mustUUID(v)
	}
	if v, ok := values["name"].(string); ok {
		product.Name = v
	}
	if v, ok := values["description"].(string); ok {
		product.Description = v
	}
	if v, ok := values["price"]; ok {
		product.Price = mustFloat(v)
	}
	if v, ok := values["image"].(string); ok {
		product.Image = v
	}
	if v, ok := values["status"].(string); ok {
		product.Status = v
	}
	if v, ok := values["sort_order"]; ok {
		product.SortOrder = mustInt(v)
	}
}

func orderFromColumns(columns []string, args []driver.Value) model.Order {
	values := mapColumns(columns, args)
	order := model.Order{ID: uuid.New()}
	if v, ok := values["id"]; ok {
		order.ID = mustUUID(v)
	}
	if v, ok := values["user_id"]; ok {
		order.UserID = mustUUID(v)
	}
	order.OrderNo, _ = values["order_no"].(string)
	order.TotalAmount = mustFloat(values["total_amount"])
	order.Status, _ = values["status"].(string)
	order.Remark, _ = values["remark"].(string)
	return order
}

func orderItemFromColumns(columns []string, args []driver.Value) model.OrderItem {
	values := mapColumns(columns, args)
	item := model.OrderItem{ID: uuid.New()}
	if v, ok := values["id"]; ok {
		item.ID = mustUUID(v)
	}
	if v, ok := values["order_id"]; ok {
		item.OrderID = mustUUID(v)
	}
	if v, ok := values["product_id"]; ok {
		item.ProductID = mustUUID(v)
	}
	item.Quantity = mustInt(values["quantity"])
	item.Price = mustFloat(values["price"])
	return item
}

func userColumns() []string {
	return []string{"id", "username", "password", "name", "phone", "role", "created_at", "updated_at", "deleted_at"}
}
func userValues(user model.User) []driver.Value {
	return []driver.Value{user.ID.String(), user.Username, user.Password, user.Name, user.Phone, user.Role, time.Time{}, time.Time{}, nil}
}

func categoryColumns() []string {
	return []string{"id", "name", "sort_order", "created_at", "updated_at", "deleted_at"}
}
func categoryValues(category model.Category) []driver.Value {
	return []driver.Value{category.ID.String(), category.Name, category.SortOrder, time.Time{}, time.Time{}, nil}
}

func productColumns() []string {
	return []string{"id", "category_id", "name", "description", "price", "image", "status", "sort_order", "created_at", "updated_at", "deleted_at"}
}
func productValues(product model.Product) []driver.Value {
	return []driver.Value{product.ID.String(), product.CategoryID.String(), product.Name, product.Description, product.Price, product.Image, product.Status, product.SortOrder, time.Time{}, time.Time{}, nil}
}

func orderColumns() []string {
	return []string{"id", "user_id", "order_no", "total_amount", "status", "remark", "created_at", "updated_at", "deleted_at"}
}
func orderValues(order model.Order) []driver.Value {
	return []driver.Value{order.ID.String(), order.UserID.String(), order.OrderNo, order.TotalAmount, order.Status, order.Remark, time.Time{}, time.Time{}, nil}
}

func orderItemColumns() []string {
	return []string{"id", "order_id", "product_id", "quantity", "price"}
}
func orderItemValues(item model.OrderItem) []driver.Value {
	return []driver.Value{item.ID.String(), item.OrderID.String(), item.ProductID.String(), item.Quantity, item.Price}
}

func matchesUser(q string, user model.User, args []driver.Value) bool {
	switch {
	case strings.Contains(q, "username ="):
		return user.Username == args[0]
	case strings.Contains(q, "phone ="):
		return user.Phone == args[0]
	case strings.Contains(q, "id ="):
		return user.ID == mustUUID(args[0])
	default:
		return true
	}
}

func matchesProduct(q string, product model.Product, args []driver.Value) bool {
	if strings.Contains(q, " id =") || strings.Contains(q, "where id =") {
		return product.ID == mustUUID(args[0])
	}
	arg := 0
	if strings.Contains(q, "category_id =") {
		if product.CategoryID != mustUUID(args[arg]) {
			return false
		}
		arg++
	}
	if strings.Contains(q, "status =") {
		return product.Status == args[arg]
	}
	if strings.Contains(q, " id in") {
		for _, candidate := range args {
			if product.ID == mustUUID(candidate) {
				return true
			}
		}
		return false
	}
	return true
}

func matchesOrder(q string, order model.Order, args []driver.Value) bool {
	switch {
	case strings.Contains(q, "user_id ="):
		return order.UserID == mustUUID(args[0])
	case strings.Contains(q, "id ="):
		return order.ID == mustUUID(args[0])
	default:
		return true
	}
}

func matchesOrderItem(q string, item model.OrderItem, args []driver.Value) bool {
	if strings.Contains(q, "order_id =") || strings.Contains(q, "order_id in") {
		for _, candidate := range args {
			if item.OrderID == mustUUID(candidate) {
				return true
			}
		}
		return false
	}
	return true
}

func sortedProducts(products map[uuid.UUID]model.Product) []model.Product {
	result := make([]model.Product, 0, len(products))
	for _, product := range products {
		result = append(result, product)
	}
	sort.Slice(result, func(i, j int) bool { return result[i].Name < result[j].Name })
	return result
}

func sortedOrders(orders map[uuid.UUID]model.Order) []model.Order {
	result := make([]model.Order, 0, len(orders))
	for _, order := range orders {
		result = append(result, order)
	}
	sort.Slice(result, func(i, j int) bool { return result[i].OrderNo < result[j].OrderNo })
	return result
}

func sortedOrderItems(items map[uuid.UUID]model.OrderItem) []model.OrderItem {
	result := make([]model.OrderItem, 0, len(items))
	for _, item := range items {
		result = append(result, item)
	}
	sort.Slice(result, func(i, j int) bool { return result[i].ID.String() < result[j].ID.String() })
	return result
}

func mustUUID(value driver.Value) uuid.UUID {
	switch v := value.(type) {
	case uuid.UUID:
		return v
	case string:
		id, err := uuid.Parse(v)
		if err != nil {
			panic(err)
		}
		return id
	case []byte:
		id, err := uuid.ParseBytes(v)
		if err != nil {
			panic(err)
		}
		return id
	default:
		panic(fmt.Sprintf("unsupported uuid value %T: %v", value, value))
	}
}

func mustFloat(value driver.Value) float64 {
	switch v := value.(type) {
	case nil:
		return 0
	case float64:
		return v
	case int64:
		return float64(v)
	case string:
		var f float64
		_, _ = fmt.Sscanf(v, "%f", &f)
		return f
	default:
		panic(fmt.Sprintf("unsupported float value %T: %v", value, value))
	}
}

func mustInt(value driver.Value) int {
	switch v := value.(type) {
	case nil:
		return 0
	case int64:
		return int(v)
	case int:
		return v
	default:
		panic(fmt.Sprintf("unsupported int value %T: %v", value, value))
	}
}

func TestAuthServiceRegister(t *testing.T) {
	existingPassword, err := bcrypt.GenerateFromPassword([]byte("secret"), bcrypt.DefaultCost)
	if err != nil {
		t.Fatalf("hash password: %v", err)
	}

	tests := []struct {
		name    string
		seed    []model.User
		input   model.User
		wantErr string
	}{
		{
			name:  "creates customer with hashed password",
			input: model.User{Username: "alice", Password: "secret", Name: "Alice", Phone: "13800000000"},
		},
		{
			name:    "rejects duplicate username",
			seed:    []model.User{{ID: uuid.New(), Username: "alice", Password: string(existingPassword), Phone: "13800000000", Role: "customer"}},
			input:   model.User{Username: "alice", Password: "secret", Name: "Alice Two", Phone: "13900000000"},
			wantErr: "username already exists",
		},
		{
			name:    "rejects duplicate phone",
			seed:    []model.User{{ID: uuid.New(), Username: "bob", Password: string(existingPassword), Phone: "13800000000", Role: "customer"}},
			input:   model.User{Username: "alice", Password: "secret", Name: "Alice", Phone: "13800000000"},
			wantErr: "phone already exists",
		},
		{
			name:  "allows empty phone boundary",
			input: model.User{Username: "alice", Password: "secret", Name: "Alice"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			store := newServiceTestStore()
			for _, user := range tt.seed {
				store.users[user.ID] = user
			}
			userRepo, _, _ := newTestRepositories(t, store)
			service := NewAuthService(config.JWTConfig{Secret: "test-secret", Expiration: 1}, userRepo)

			user, err := service.Register(tt.input.Username, tt.input.Password, tt.input.Name, tt.input.Phone)
			if tt.wantErr != "" {
				if err == nil || err.Error() != tt.wantErr {
					t.Fatalf("Register() error = %v, want %q", err, tt.wantErr)
				}
				return
			}

			if err != nil {
				t.Fatalf("Register() error = %v", err)
			}
			if user == nil || user.Username != tt.input.Username || user.Role != "customer" {
				t.Fatalf("Register() user = %#v", user)
			}
			if user.Password == tt.input.Password {
				t.Fatal("Register() stored plaintext password")
			}
			if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(tt.input.Password)); err != nil {
				t.Fatalf("Register() password hash mismatch: %v", err)
			}
		})
	}
}

func TestAuthServiceLogin(t *testing.T) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte("secret"), bcrypt.DefaultCost)
	if err != nil {
		t.Fatalf("hash password: %v", err)
	}
	userID := uuid.New()

	tests := []struct {
		name     string
		username string
		password string
		seed     []model.User
		wantErr  string
	}{
		{
			name:     "returns signed token",
			username: "alice",
			password: "secret",
			seed:     []model.User{{ID: userID, Username: "alice", Password: string(hashedPassword), Role: "admin"}},
		},
		{
			name:     "rejects unknown user",
			username: "missing",
			password: "secret",
			wantErr:  "invalid credentials",
		},
		{
			name:     "rejects wrong password",
			username: "alice",
			password: "wrong",
			seed:     []model.User{{ID: userID, Username: "alice", Password: string(hashedPassword), Role: "customer"}},
			wantErr:  "invalid credentials",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			store := newServiceTestStore()
			for _, user := range tt.seed {
				store.users[user.ID] = user
			}
			userRepo, _, _ := newTestRepositories(t, store)
			service := NewAuthService(config.JWTConfig{Secret: "test-secret", Expiration: 1}, userRepo)

			tokenString, err := service.Login(tt.username, tt.password)
			if tt.wantErr != "" {
				if err == nil || err.Error() != tt.wantErr {
					t.Fatalf("Login() error = %v, want %q", err, tt.wantErr)
				}
				return
			}
			if err != nil {
				t.Fatalf("Login() error = %v", err)
			}
			claims := jwt.MapClaims{}
			token, err := jwt.ParseWithClaims(tokenString, claims, func(*jwt.Token) (interface{}, error) {
				return []byte("test-secret"), nil
			})
			if err != nil || !token.Valid {
				t.Fatalf("Login() token invalid: token=%v err=%v", token, err)
			}
			if claims["user_id"] != userID.String() || claims["username"] != tt.username {
				t.Fatalf("Login() claims = %#v", claims)
			}
		})
	}
}

func TestAuthServiceGetMeTokenClaims(t *testing.T) {
	userID := uuid.New()
	service := NewAuthService(config.JWTConfig{Secret: "test-secret", Expiration: 1}, nil)

	tests := []struct {
		name       string
		makeToken  func(t *testing.T) string
		wantUserID uuid.UUID
		wantName   string
		wantRole   string
		wantErr    bool
	}{
		{
			name: "valid token returns current user claims",
			makeToken: func(t *testing.T) string {
				t.Helper()
				return signTestToken(t, "test-secret", jwt.MapClaims{"user_id": userID.String(), "username": "alice", "role": "admin", "exp": time.Now().Add(time.Hour).Unix()})
			},
			wantUserID: userID,
			wantName:   "alice",
			wantRole:   "admin",
		},
		{
			name: "rejects malformed user id",
			makeToken: func(t *testing.T) string {
				t.Helper()
				return signTestToken(t, "test-secret", jwt.MapClaims{"user_id": "bad", "username": "alice", "role": "customer", "exp": time.Now().Add(time.Hour).Unix()})
			},
			wantErr: true,
		},
		{
			name: "rejects wrong secret",
			makeToken: func(t *testing.T) string {
				t.Helper()
				return signTestToken(t, "wrong-secret", jwt.MapClaims{"user_id": userID.String(), "username": "alice", "role": "customer", "exp": time.Now().Add(time.Hour).Unix()})
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotUserID, gotUsername, gotRole, err := service.ValidateToken(tt.makeToken(t))
			if tt.wantErr {
				if err == nil {
					t.Fatal("ValidateToken() expected error")
				}
				return
			}
			if err != nil {
				t.Fatalf("ValidateToken() error = %v", err)
			}
			if gotUserID != tt.wantUserID || gotUsername != tt.wantName || gotRole != tt.wantRole {
				t.Fatalf("ValidateToken() = (%s, %s, %s)", gotUserID, gotUsername, gotRole)
			}
		})
	}
}

func signTestToken(t *testing.T, secret string, claims jwt.MapClaims) string {
	t.Helper()
	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(secret))
	if err != nil {
		t.Fatalf("sign test token: %v", err)
	}
	return token
}

func assertEqual[T comparable](t *testing.T, name string, got, want T) {
	t.Helper()
	if got != want {
		t.Fatalf("%s = %v, want %v", name, got, want)
	}
}

func assertDeepEqual(t *testing.T, name string, got, want interface{}) {
	t.Helper()
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("%s = %#v, want %#v", name, got, want)
	}
}
