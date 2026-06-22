<script lang="ts">
	import { onMount } from 'svelte';
	import { products } from '$lib/stores';
	import { cart } from '$lib/stores';
	import { productsApi } from '$lib/api';
	import type { Product } from '$lib/api';
	import { Button } from '$lib/components/ui/button';
	import { Plus, Coffee, Search } from 'lucide-svelte';
	import { cn } from '$lib/utils';

	const categories = ['All', 'Coffee', 'Tea', 'Pastry', 'Special'];

	let selectedCategory = $state('All');
	let searchQuery = $state('');

	let filteredProducts = $derived(
		$products.products.filter((p) => {
			const matchesCategory =
				selectedCategory === 'All' || p.category?.name === selectedCategory;
			const matchesSearch =
				!searchQuery ||
				p.name.toLowerCase().includes(searchQuery.toLowerCase());
			return matchesCategory && matchesSearch;
		})
	);

	onMount(async () => {
		products.setLoading();
		try {
			const data = await productsApi.list();
			products.setProducts(data);
		} catch (err) {
			products.setError(err instanceof Error ? err.message : 'Failed to load products');
		}
	});

	function addToCart(product: Product) {
		cart.add(product, 1);
	}
</script>

<div class="mx-auto max-w-lg">
	<header class="sticky top-0 z-10 bg-background/95 backdrop-blur supports-[backdrop-filter]:bg-background/60">
		<div class="px-4 pt-4 pb-2">
			<div class="flex items-center gap-2">
				<Coffee class="h-6 w-6 text-primary" />
				<h1 class="text-xl font-bold text-foreground">Coffee Shop</h1>
			</div>
		</div>

		<div class="px-4 pb-3">
			<div class="relative">
				<Search class="absolute left-3 top-1/2 h-4 w-4 -translate-y-1/2 text-muted-foreground" />
				<input
					type="text"
					placeholder="Search menu..."
					bind:value={searchQuery}
					class="h-10 w-full rounded-lg border border-input bg-background pl-9 pr-4 text-sm placeholder:text-muted-foreground focus:outline-none focus:ring-2 focus:ring-ring"
				/>
			</div>
		</div>

		<div class="flex gap-2 overflow-x-auto px-4 pb-3 scrollbar-none">
			{#each categories as category (category)}
				<button
					onclick={() => (selectedCategory = category)}
					class={cn(
						'whitespace-nowrap rounded-full px-4 py-1.5 text-sm font-medium transition-colors',
						selectedCategory === category
							? 'bg-primary text-primary-foreground'
							: 'bg-secondary text-secondary-foreground hover:bg-secondary/80'
					)}
				>
					{category}
				</button>
			{/each}
		</div>
	</header>

	{#if $products.loading}
		<div class="flex flex-col items-center justify-center px-4 py-16 text-center">
			<div class="h-8 w-8 animate-spin rounded-full border-4 border-primary border-t-transparent"></div>
			<p class="mt-4 text-sm text-muted-foreground">Loading menu...</p>
		</div>
	{:else if $products.error}
		<div class="flex flex-col items-center justify-center px-4 py-16 text-center">
			<p class="text-sm text-destructive">{$products.error}</p>
		</div>
	{:else}
		<div class="grid grid-cols-2 gap-3 px-4 sm:grid-cols-3">
			{#each filteredProducts as product (product.id)}
				<div class="group overflow-hidden rounded-lg border border-border bg-card transition-shadow hover:shadow-md">
					<div class="aspect-square bg-muted">
						{#if product.image}
							<img
								src={product.image}
								alt={product.name}
								class="h-full w-full object-cover"
							/>
						{:else}
							<div class="flex h-full items-center justify-center">
								<Coffee class="h-10 w-10 text-muted-foreground/50" />
							</div>
						{/if}
					</div>
					<div class="p-3">
						<h3 class="truncate text-sm font-semibold text-card-foreground">
							{product.name}
						</h3>
						<p class="mt-0.5 text-xs text-muted-foreground">{product.category?.name ?? ''}</p>
						<div class="mt-2 flex items-center justify-between">
							<span class="text-sm font-bold text-foreground">
								${product.price.toFixed(2)}
							</span>
							<Button
								size="icon"
								variant="ghost"
								class="h-8 w-8"
								onclick={() => addToCart(product)}
							>
								<Plus class="h-4 w-4" />
							</Button>
						</div>
					</div>
				</div>
			{/each}
		</div>

		{#if filteredProducts.length === 0}
			<div class="flex flex-col items-center justify-center px-4 py-16 text-center">
				<Coffee class="h-12 w-12 text-muted-foreground/50" />
				<p class="mt-4 text-sm text-muted-foreground">No items found</p>
			</div>
		{/if}
	{/if}
</div>
