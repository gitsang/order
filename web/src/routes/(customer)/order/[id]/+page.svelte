<script lang="ts">
	import { onMount } from 'svelte';
	import { page } from '$app/stores';
	import { ordersApi } from '$lib/api';
	import type { Order } from '$lib/api';
	import { cn } from '$lib/utils';
	import { Button } from '$lib/components/ui/button';
	import { ArrowLeft, Package, Clock, CheckCircle, XCircle } from 'lucide-svelte';
	import { goto } from '$app/navigation';

	let order = $state<Order | null>(null);
	let loading = $state(true);
	let error = $state('');

	const orderId = $derived($page.params.id ?? '');

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
			order = await ordersApi.get(orderId);
		} catch (err) {
			error = err instanceof Error ? err.message : 'Failed to load order';
		} finally {
			loading = false;
		}
	});

	function formatDate(dateStr: string | undefined) {
		if (!dateStr) return '';
		return new Date(dateStr).toLocaleDateString('en-US', {
			weekday: 'short',
			month: 'short',
			day: 'numeric',
			year: 'numeric',
			hour: '2-digit',
			minute: '2-digit'
		});
	}
</script>

<div class="mx-auto max-w-lg">
	<header class="sticky top-0 z-10 flex items-center gap-3 bg-background/95 px-4 py-3 backdrop-blur supports-[backdrop-filter]:bg-background/60">
		<Button variant="ghost" size="icon" class="h-9 w-9" onclick={() => goto('/order')}>
			<ArrowLeft class="h-5 w-5" />
		</Button>
		<h1 class="text-lg font-semibold text-foreground">Order Details</h1>
	</header>

	{#if loading}
		<div class="flex items-center justify-center px-4 py-20">
			<div class="h-8 w-8 animate-spin rounded-full border-4 border-primary border-t-transparent"></div>
		</div>
	{:else if error}
		<div class="flex flex-col items-center justify-center px-4 py-20 text-center">
			<XCircle class="h-12 w-12 text-destructive/50" />
			<p class="mt-4 text-sm text-destructive">{error}</p>
			<Button class="mt-4" variant="outline" onclick={() => goto('/order')}>
				Back to Orders
			</Button>
		</div>
	{:else if order}
		<div class="space-y-4 px-4 py-2">
			<div class="rounded-lg border border-border bg-card p-4">
				<div class="flex items-center justify-between">
					<h2 class="text-lg font-semibold text-card-foreground">
						Order #{order.order_no ?? order.id.slice(-6)}
					</h2>
					<span
						class={cn(
							'rounded-full px-3 py-1 text-xs font-medium capitalize',
							statusColors[order.status] ?? 'bg-secondary text-secondary-foreground'
						)}
					>
						{order.status}
					</span>
				</div>
				<p class="mt-1 text-sm text-muted-foreground">{formatDate(order.created_at)}</p>
			</div>

			<div class="rounded-lg border border-border bg-card">
				<div class="border-b border-border px-4 py-3">
					<h3 class="text-sm font-semibold text-card-foreground">Items</h3>
				</div>
				<div class="divide-y divide-border">
					{#each order.items ?? [] as item (item.id)}
						<div class="flex items-center gap-3 px-4 py-3">
							<div class="flex h-10 w-10 shrink-0 items-center justify-center rounded-lg bg-muted overflow-hidden">
								{#if item.product?.image}
									<img src={item.product.image} alt={item.product?.name ?? ''} class="h-full w-full object-cover" />
								{:else}
									<Package class="h-4 w-4 text-muted-foreground/50" />
								{/if}
							</div>
							<div class="min-w-0 flex-1">
								<p class="truncate text-sm font-medium text-foreground">
									{item.product?.name ?? `Product #${item.product_id.slice(-6)}`}
								</p>
								<p class="text-xs text-muted-foreground">Qty: {item.quantity}</p>
							</div>
							<span class="text-sm font-medium text-foreground">
								${(item.price * item.quantity).toFixed(2)}
							</span>
						</div>
					{/each}
				</div>
			</div>

			<div class="rounded-lg border border-border bg-card p-4">
				<div class="flex items-center justify-between">
					<span class="text-sm text-muted-foreground">Total</span>
					<span class="text-lg font-bold text-foreground">${order.total_amount.toFixed(2)}</span>
				</div>
			</div>

			{#if order.remark}
				<div class="rounded-lg border border-border bg-card p-4">
					<h3 class="text-sm font-semibold text-card-foreground">Remark</h3>
					<p class="mt-1 text-sm text-muted-foreground">{order.remark}</p>
				</div>
			{/if}

			<div class="py-4">
				<Button class="w-full" variant="outline" onclick={() => goto('/')}>
					Order Again
				</Button>
			</div>
		</div>
	{/if}
</div>
