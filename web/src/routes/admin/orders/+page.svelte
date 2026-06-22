<script lang="ts">
	type OrderStatus = 'pending' | 'preparing' | 'ready' | 'completed' | 'cancelled';

	let statusFilter = $state<OrderStatus | 'all'>('all');
	let searchQuery = $state('');

	let orders = $state([
		{ id: 'ORD-001', customer: 'Alice Chen', items: ['Espresso', 'Latte', 'Croissant'], total: 11.0, status: 'pending' as OrderStatus, createdAt: '2024-01-15 09:30' },
		{ id: 'ORD-002', customer: 'Bob Kim', items: ['Cappuccino'], total: 4.0, status: 'preparing' as OrderStatus, createdAt: '2024-01-15 09:45' },
		{ id: 'ORD-003', customer: 'Carol Liu', items: ['Mocha', 'Muffin', 'Espresso', 'Latte', 'Tea'], total: 21.5, status: 'ready' as OrderStatus, createdAt: '2024-01-15 10:00' },
		{ id: 'ORD-004', customer: 'David Park', items: ['Americano', 'Scone'], total: 7.5, status: 'completed' as OrderStatus, createdAt: '2024-01-15 08:15' },
		{ id: 'ORD-005', customer: 'Eve Zhang', items: ['Latte'], total: 4.5, status: 'completed' as OrderStatus, createdAt: '2024-01-15 08:30' },
		{ id: 'ORD-006', customer: 'Frank Wu', items: ['Espresso', 'Muffin'], total: 6.5, status: 'cancelled' as OrderStatus, createdAt: '2024-01-15 07:45' },
		{ id: 'ORD-007', customer: 'Grace Li', items: ['Cappuccino', 'Croissant', 'Juice'], total: 12.0, status: 'pending' as OrderStatus, createdAt: '2024-01-15 10:15' },
		{ id: 'ORD-008', customer: 'Henry Zhao', items: ['Mocha', 'Latte'], total: 9.5, status: 'preparing' as OrderStatus, createdAt: '2024-01-15 10:30' }
	]);

	const statusTabs: { value: OrderStatus | 'all'; label: string }[] = [
		{ value: 'all', label: 'All' },
		{ value: 'pending', label: 'Pending' },
		{ value: 'preparing', label: 'Preparing' },
		{ value: 'ready', label: 'Ready' },
		{ value: 'completed', label: 'Completed' },
		{ value: 'cancelled', label: 'Cancelled' }
	];

	const statusColors: Record<OrderStatus, string> = {
		pending: 'bg-yellow-100 text-yellow-800',
		preparing: 'bg-blue-100 text-blue-800',
		ready: 'bg-green-100 text-green-800',
		completed: 'bg-muted text-muted-foreground',
		cancelled: 'bg-red-100 text-red-800'
	};

	const nextStatus: Record<Exclude<OrderStatus, 'completed' | 'cancelled'>, OrderStatus> = {
		pending: 'preparing',
		preparing: 'ready',
		ready: 'completed'
	};

	let filteredOrders = $derived(
		orders.filter((o) => {
			const matchesStatus = statusFilter === 'all' || o.status === statusFilter;
			const matchesSearch =
				o.id.toLowerCase().includes(searchQuery.toLowerCase()) ||
				o.customer.toLowerCase().includes(searchQuery.toLowerCase());
			return matchesStatus && matchesSearch;
		})
	);

	let statusCounts = $derived(
		orders.reduce(
			(acc, o) => {
				acc[o.status] = (acc[o.status] || 0) + 1;
				return acc;
			},
			{} as Record<string, number>
		)
	);

	function advanceStatus(id: string) {
		orders = orders.map((o) => {
			if (o.id !== id) return o;
			if (o.status === 'completed' || o.status === 'cancelled') return o;
			return { ...o, status: nextStatus[o.status] };
		});
	}

	function cancelOrder(id: string) {
		orders = orders.map((o) => (o.id === id ? { ...o, status: 'cancelled' as OrderStatus } : o));
	}
</script>

<svelte:head>
	<title>Orders - Coffee Admin</title>
</svelte:head>

<div class="space-y-6">
	<div>
		<h1 class="text-2xl font-bold text-foreground">Orders</h1>
		<p class="mt-1 text-muted-foreground">{orders.length} orders total</p>
	</div>

	<div class="flex flex-col gap-4 sm:flex-row sm:items-center sm:justify-between">
		<div class="flex flex-wrap gap-2">
			{#each statusTabs as tab}
				{@const count = tab.value === 'all' ? orders.length : (statusCounts[tab.value] || 0)}
				<button
					class="inline-flex items-center gap-1.5 rounded-md px-3 py-1.5 text-sm font-medium transition-colors {statusFilter === tab.value
						? 'bg-primary text-primary-foreground'
						: 'border border-border text-muted-foreground hover:bg-accent hover:text-accent-foreground'}"
					onclick={() => (statusFilter = tab.value)}
				>
					{tab.label}
					<span class="rounded-full bg-background/20 px-1.5 py-0.5 text-xs">{count}</span>
				</button>
			{/each}
		</div>

		<div class="relative w-full sm:w-64">
			<svg class="absolute left-3 top-1/2 h-4 w-4 -translate-y-1/2 text-muted-foreground" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
				<circle cx="11" cy="11" r="8" />
				<path d="m21 21-4.3-4.3" />
			</svg>
			<input
				type="text"
				bind:value={searchQuery}
				placeholder="Search orders..."
				class="w-full rounded-md border border-border bg-background py-2 pl-10 pr-3 text-sm placeholder:text-muted-foreground focus:outline-none focus:ring-2 focus:ring-ring"
			/>
		</div>
	</div>

	<div class="rounded-lg border border-border bg-card">
		<div class="overflow-x-auto">
			<table class="w-full">
				<thead>
					<tr class="border-b border-border bg-muted/50">
						<th class="px-6 py-3 text-left text-xs font-medium uppercase tracking-wider text-muted-foreground">Order ID</th>
						<th class="px-6 py-3 text-left text-xs font-medium uppercase tracking-wider text-muted-foreground">Customer</th>
						<th class="px-6 py-3 text-left text-xs font-medium uppercase tracking-wider text-muted-foreground">Items</th>
						<th class="px-6 py-3 text-left text-xs font-medium uppercase tracking-wider text-muted-foreground">Total</th>
						<th class="px-6 py-3 text-left text-xs font-medium uppercase tracking-wider text-muted-foreground">Status</th>
						<th class="px-6 py-3 text-left text-xs font-medium uppercase tracking-wider text-muted-foreground">Time</th>
						<th class="px-6 py-3 text-right text-xs font-medium uppercase tracking-wider text-muted-foreground">Actions</th>
					</tr>
				</thead>
				<tbody class="divide-y divide-border">
					{#each filteredOrders as order (order.id)}
						<tr class="hover:bg-muted/50">
							<td class="whitespace-nowrap px-6 py-4 text-sm font-medium text-foreground">{order.id}</td>
							<td class="whitespace-nowrap px-6 py-4 text-sm text-muted-foreground">{order.customer}</td>
							<td class="px-6 py-4 text-sm text-muted-foreground">
								<div class="flex flex-wrap gap-1">
									{#each order.items.slice(0, 3) as item}
										<span class="inline-flex rounded bg-muted px-1.5 py-0.5 text-xs">{item}</span>
									{/each}
									{#if order.items.length > 3}
										<span class="inline-flex rounded bg-muted px-1.5 py-0.5 text-xs text-muted-foreground">
											+{order.items.length - 3}
										</span>
									{/if}
								</div>
							</td>
							<td class="whitespace-nowrap px-6 py-4 text-sm font-medium text-foreground">${order.total.toFixed(2)}</td>
							<td class="whitespace-nowrap px-6 py-4">
								<span class="inline-flex rounded-full px-2 py-1 text-xs font-medium {statusColors[order.status]}">
									{order.status}
								</span>
							</td>
							<td class="whitespace-nowrap px-6 py-4 text-sm text-muted-foreground">{order.createdAt}</td>
							<td class="whitespace-nowrap px-6 py-4 text-right">
								<div class="flex items-center justify-end gap-2">
									{#if order.status !== 'completed' && order.status !== 'cancelled'}
										<button
											class="rounded-md px-2.5 py-1 text-xs font-medium text-primary transition-colors hover:bg-primary/10"
											onclick={() => advanceStatus(order.id)}
										>
											{#if order.status === 'pending'}
												Start Preparing
											{:else if order.status === 'preparing'}
												Mark Ready
											{:else if order.status === 'ready'}
												Complete
											{/if}
										</button>
										<button
											class="rounded-md p-1.5 text-destructive transition-colors hover:bg-destructive/10"
											onclick={() => cancelOrder(order.id)}
											title="Cancel order"
										>
											<svg class="h-4 w-4" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
												<path d="M18 6 6 18M6 6l12 12" />
											</svg>
										</button>
									{:else}
										<span class="text-xs text-muted-foreground">No actions</span>
									{/if}
								</div>
							</td>
						</tr>
					{:else}
						<tr>
							<td colspan="7" class="px-6 py-12 text-center text-muted-foreground">
								No orders found
							</td>
						</tr>
					{/each}
				</tbody>
			</table>
		</div>
	</div>
</div>
