<script lang="ts">
	const stats = [
		{ label: 'Total Products', value: '24', change: '+3 this week', icon: 'package', positive: true },
		{ label: 'Active Orders', value: '12', change: '4 pending', icon: 'orders', positive: true },
		{ label: "Today's Revenue", value: '$1,248', change: '+12% vs yesterday', icon: 'revenue', positive: true },
		{ label: 'Customers', value: '156', change: '+8 this week', icon: 'customers', positive: true }
	];

	const recentOrders = [
		{ id: 'ORD-001', customer: 'Alice Chen', items: 3, total: '$24.50', status: 'preparing' },
		{ id: 'ORD-002', customer: 'Bob Kim', items: 1, total: '$8.00', status: 'ready' },
		{ id: 'ORD-003', customer: 'Carol Liu', items: 5, total: '$42.75', status: 'pending' },
		{ id: 'ORD-004', customer: 'David Park', items: 2, total: '$15.00', status: 'completed' },
		{ id: 'ORD-005', customer: 'Eve Zhang', items: 1, total: '$6.50', status: 'completed' }
	];

	const statusColors: Record<string, string> = {
		pending: 'bg-yellow-100 text-yellow-800',
		preparing: 'bg-blue-100 text-blue-800',
		ready: 'bg-green-100 text-green-800',
		completed: 'bg-muted text-muted-foreground'
	};
</script>

<svelte:head>
	<title>Dashboard - Coffee Admin</title>
</svelte:head>

<div class="space-y-8">
	<div>
		<h1 class="text-2xl font-bold text-foreground">Dashboard</h1>
		<p class="mt-1 text-muted-foreground">Overview of your coffee shop</p>
	</div>

	<div class="grid grid-cols-1 gap-6 sm:grid-cols-2 lg:grid-cols-4">
		{#each stats as stat}
			<div class="rounded-lg border border-border bg-card p-6">
				<div class="flex items-center justify-between">
					<p class="text-sm font-medium text-muted-foreground">{stat.label}</p>
					{#if stat.icon === 'package'}
						<svg class="h-5 w-5 text-muted-foreground" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
							<path d="M20 7l-8-4-8 4m16 0l-8 4m8-4v10l-8 4m0-10L4 7m8 4v10M4 7v10l8 4" />
						</svg>
					{:else if stat.icon === 'orders'}
						<svg class="h-5 w-5 text-muted-foreground" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
							<path d="M16 11V7a4 4 0 00-8 0v4M5 9h14l1 12H4L5 9z" />
						</svg>
					{:else if stat.icon === 'revenue'}
						<svg class="h-5 w-5 text-muted-foreground" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
							<path d="M12 8c-1.657 0-3 .895-3 2s1.343 2 3 2 3 .895 3 2-1.343 2-3 2m0-8c1.11 0 2.08.402 2.599 1M12 8V7m0 1v8m0 0v1m0-1c-1.11 0-2.08-.402-2.599-1M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
						</svg>
					{:else if stat.icon === 'customers'}
						<svg class="h-5 w-5 text-muted-foreground" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
							<path d="M17 20h5v-2a3 3 0 00-5.356-1.857M17 20H7m10 0v-2c0-.656-.126-1.283-.356-1.857M7 20H2v-2a3 3 0 015.356-1.857M7 20v-2c0-.656.126-1.283.356-1.857m0 0a5.002 5.002 0 019.288 0M15 7a3 3 0 11-6 0 3 3 0 016 0zm6 3a2 2 0 11-4 0 2 2 0 014 0zM7 10a2 2 0 11-4 0 2 2 0 014 0z" />
						</svg>
					{/if}
				</div>
				<p class="mt-2 text-3xl font-bold text-foreground">{stat.value}</p>
				<p class="mt-1 text-sm {stat.positive ? 'text-green-600' : 'text-red-600'}">{stat.change}</p>
			</div>
		{/each}
	</div>

	<div class="rounded-lg border border-border bg-card">
		<div class="flex items-center justify-between border-b border-border px-6 py-4">
			<h2 class="text-lg font-semibold text-foreground">Recent Orders</h2>
			<a href="/admin/orders" class="text-sm font-medium text-primary underline-offset-4 hover:underline">
				View all
			</a>
		</div>
		<div class="overflow-x-auto">
			<table class="w-full">
				<thead>
					<tr class="border-b border-border">
						<th class="px-6 py-3 text-left text-xs font-medium uppercase tracking-wider text-muted-foreground">Order ID</th>
						<th class="px-6 py-3 text-left text-xs font-medium uppercase tracking-wider text-muted-foreground">Customer</th>
						<th class="px-6 py-3 text-left text-xs font-medium uppercase tracking-wider text-muted-foreground">Items</th>
						<th class="px-6 py-3 text-left text-xs font-medium uppercase tracking-wider text-muted-foreground">Total</th>
						<th class="px-6 py-3 text-left text-xs font-medium uppercase tracking-wider text-muted-foreground">Status</th>
					</tr>
				</thead>
				<tbody class="divide-y divide-border">
					{#each recentOrders as order}
						<tr class="hover:bg-muted/50">
							<td class="whitespace-nowrap px-6 py-4 text-sm font-medium text-foreground">{order.id}</td>
							<td class="whitespace-nowrap px-6 py-4 text-sm text-muted-foreground">{order.customer}</td>
							<td class="whitespace-nowrap px-6 py-4 text-sm text-muted-foreground">{order.items}</td>
							<td class="whitespace-nowrap px-6 py-4 text-sm text-muted-foreground">{order.total}</td>
							<td class="whitespace-nowrap px-6 py-4">
								<span class="inline-flex rounded-full px-2 py-1 text-xs font-medium {statusColors[order.status]}">
									{order.status}
								</span>
							</td>
						</tr>
					{/each}
				</tbody>
			</table>
		</div>
	</div>

	<div class="rounded-lg border border-border bg-card p-6">
		<h2 class="text-lg font-semibold text-foreground mb-4">Quick Actions</h2>
		<div class="flex flex-wrap gap-3">
			<a href="/admin/products" class="inline-flex items-center rounded-md bg-primary px-4 py-2 text-sm font-medium text-primary-foreground transition-colors hover:bg-primary/90">
				Manage Products
			</a>
			<a href="/admin/orders" class="inline-flex items-center rounded-md bg-primary px-4 py-2 text-sm font-medium text-primary-foreground transition-colors hover:bg-primary/90">
				View Orders
			</a>
		</div>
	</div>
</div>
