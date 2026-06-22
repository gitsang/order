<script lang="ts">
	import { cart } from '$lib/stores';
	import { ordersApi } from '$lib/api';
	import type { Order, CreateOrderRequest } from '$lib/api';
	import { Button } from '$lib/components/ui/button';
	import { Minus, Plus, Trash2, ShoppingCart, ArrowRight } from 'lucide-svelte';
	import { goto } from '$app/navigation';

	let isSubmitting = $state(false);

	let isEmpty = $derived($cart.items.length === 0);

	function increment(productId: string) {
		const item = $cart.items.find((i) => i.product.id === productId);
		if (item) {
			cart.updateQuantity(productId, item.quantity + 1);
		}
	}

	function decrement(productId: string) {
		const item = $cart.items.find((i) => i.product.id === productId);
		if (item && item.quantity > 1) {
			cart.updateQuantity(productId, item.quantity - 1);
		} else {
			cart.remove(productId);
		}
	}

	async function handleCheckout() {
		if (isEmpty) return;
		isSubmitting = true;
		try {
			const order = await ordersApi.create({
				items: $cart.items.map((item) => ({
					product_id: item.product.id,
					quantity: item.quantity
				}))
			});
			cart.clear();
			goto(`/order/${order.id}`);
		} catch (err) {
			console.error('Checkout failed:', err);
		} finally {
			isSubmitting = false;
		}
	}
</script>

<div class="mx-auto max-w-lg">
	<header class="px-4 pt-4 pb-2">
		<h1 class="text-xl font-bold text-foreground">Your Cart</h1>
		{#if !isEmpty}
			<p class="text-sm text-muted-foreground">
				{$cart.totalItems} item{$cart.totalItems !== 1 ? 's' : ''}
			</p>
		{/if}
	</header>

	{#if isEmpty}
		<div class="flex flex-col items-center justify-center px-4 py-20 text-center">
			<ShoppingCart class="h-16 w-16 text-muted-foreground/30" />
			<p class="mt-4 text-lg font-medium text-foreground">Your cart is empty</p>
			<p class="mt-1 text-sm text-muted-foreground">Add some delicious coffee to get started</p>
			<Button class="mt-6" onclick={() => goto('/')}>Browse Menu</Button>
		</div>
	{:else}
		<div class="divide-y divide-border px-4">
			{#each $cart.items as item (item.product.id)}
				<div class="flex items-center gap-3 py-4">
					<div class="flex h-14 w-14 shrink-0 items-center justify-center rounded-lg bg-muted overflow-hidden">
						{#if item.product.image}
							<img src={item.product.image} alt={item.product.name} class="h-full w-full object-cover" />
						{:else}
							<ShoppingCart class="h-6 w-6 text-muted-foreground/50" />
						{/if}
					</div>
					<div class="min-w-0 flex-1">
						<h3 class="truncate text-sm font-semibold text-foreground">{item.product.name}</h3>
						<p class="text-sm text-muted-foreground">${item.product.price.toFixed(2)}</p>
					</div>
					<div class="flex items-center gap-1">
						<Button
							variant="outline"
							size="icon"
							class="h-8 w-8"
							onclick={() => decrement(item.product.id)}
						>
							<Minus class="h-3 w-3" />
						</Button>
						<span class="w-8 text-center text-sm font-medium">{item.quantity}</span>
						<Button
							variant="outline"
							size="icon"
							class="h-8 w-8"
							onclick={() => increment(item.product.id)}
						>
							<Plus class="h-3 w-3" />
						</Button>
						<Button
							variant="ghost"
							size="icon"
							class="ml-1 h-8 w-8 text-destructive"
							onclick={() => cart.remove(item.product.id)}
						>
							<Trash2 class="h-3.5 w-3.5" />
						</Button>
					</div>
				</div>
			{/each}
		</div>

		<div class="sticky bottom-20 border-t border-border bg-background px-4 py-4">
			<div class="flex items-center justify-between">
				<span class="text-sm text-muted-foreground">Total</span>
				<span class="text-lg font-bold text-foreground">${$cart.totalAmount.toFixed(2)}</span>
			</div>
			<Button
				class="mt-3 w-full"
				size="lg"
				disabled={isSubmitting}
				onclick={handleCheckout}
			>
				{#if isSubmitting}
					Placing Order...
				{:else}
					Checkout
					<ArrowRight class="ml-2 h-4 w-4" />
				{/if}
			</Button>
		</div>
	{/if}
</div>
