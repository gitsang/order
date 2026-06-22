<script lang="ts">
	import { onMount } from 'svelte';
	import { ordersApi } from '$lib/api';
	import type { Order } from '$lib/api';
	import { cn } from '$lib/utils';
	import { ClipboardList, ChevronRight, Package } from 'lucide-svelte';

	let orders = $state<Order[]>([]);
	let loading = $state(true);
	let error = $state<string | null>(null);

	const statusColors: Record<string, string> = {
		pending: 'bg-yellow-100 text-yellow-800',
		confirmed: 'bg-blue-100 text-blue-800',
		preparing: 'bg-orange-100 text-orange-800',
		ready: 'bg-green-100 text-green-800',
		completed: 'bg-green-100 text-green-800',
		cancelled: 'bg-red-100 text-red-800'
	};

	onMount(async () => {
		try {
			orders = await ordersApi.list();
		} catch (err) {
			error = err instanceof Error ? err.message : 'Failed to load orders';
		} finally {
			loading = false;
		}
	});

	function formatDate(dateStr: string | undefined) {
		if (!dateStr) return '';
		const date = new Date(dateStr);
		return date.toLocaleDateString('en-US', {
			month: 'short',
			day: 'numeric',
			hour: '2-digit',
			minute: '2-digit'
		});
	}
</script>

<div class="mx-auto max-w-lg">
	<header class="px-4 pt-4 pb-2">
		<h1 class="text-xl font-bold text-foreground">Your Orders</h1>
	</header>

	{#if loading}
		<div class="flex flex-col items-center justify-center px-4 py-20 text-center">
			<div class="h-8 w-8 animate-spin rounded-full border-4 border-primary border-t-transparent"></div>
			<p class="mt-4 text-sm text-muted-foreground">Loading orders...</p>
		</div>
	{:else if error}
		<div class="flex flex-col items-center justify-center px-4 py-20 text-center">
			<p class="text-sm text-destructive">{error}</p>
		</div>
	{:else if orders.length === 0}
		<div class="flex flex-col items-center justify-center px-4 py-20 text-center">
			<ClipboardList class="h-16 w-16 text-muted-foreground/30" />
			<p class="mt-4 text-lg font-medium text-foreground">No orders yet</p>
			<p class="mt-1 text-sm text-muted-foreground">Your order history will appear here</p>
		</div>
	{:else}
		<div class="divide-y divide-border px-4">
			{#each orders as order (order.id)}
				<a
					href="/order/{order.id}"
					class="flex items-center gap-3 py-4 transition-colors hover:bg-muted/50"
				>
					<div class="flex h-12 w-12 shrink-0 items-center justify-center rounded-lg bg-muted">
						<Package class="h-6 w-6 text-muted-foreground" />
					</div>
					<div class="min-w-0 flex-1">
						<div class="flex items-center gap-2">
							<h3 class="text-sm font-semibold text-foreground">
								Order #{order.order_no ?? order.id.slice(-6)}
							</h3>
							<span
								class={cn(
									'rounded-full px-2 py-0.5 text-[10px] font-medium capitalize',
									statusColors[order.status] ?? 'bg-secondary text-secondary-foreground'
								)}
							>
								{order.status}
							</span>
						</div>
						<p class="mt-0.5 text-xs text-muted-foreground">
							{formatDate(order.created_at)} · {order.items?.length ?? 0} item{(order.items?.length ?? 0) !== 1 ? 's' : ''}
						</p>
						<p class="mt-0.5 text-sm font-medium text-foreground">
							${order.total_amount.toFixed(2)}
						</p>
					</div>
					<ChevronRight class="h-5 w-5 shrink-0 text-muted-foreground" />
				</a>
			{/each}
		</div>
	{/if}
</div>
