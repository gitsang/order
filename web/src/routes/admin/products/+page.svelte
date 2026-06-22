<script lang="ts">
	import { Button } from '$lib/components/ui/button';

	let showAddModal = $state(false);
	let editingId = $state<string | null>(null);
	let searchQuery = $state('');

	let products = $state([
		{ id: '1', name: 'Espresso', category: 'Coffee', price: 3.5, stock: 100, status: 'active' },
		{ id: '2', name: 'Latte', category: 'Coffee', price: 4.5, stock: 80, status: 'active' },
		{ id: '3', name: 'Cappuccino', category: 'Coffee', price: 4.0, stock: 60, status: 'active' },
		{ id: '4', name: 'Croissant', category: 'Pastry', price: 3.0, stock: 25, status: 'active' },
		{ id: '5', name: 'Mocha', category: 'Coffee', price: 5.0, stock: 5, status: 'active' },
		{ id: '6', name: 'Blueberry Muffin', category: 'Pastry', price: 3.5, stock: 0, status: 'inactive' }
	]);

	let formData = $state({ name: '', category: '', price: '', stock: '' });

	let filteredProducts = $derived(
		products.filter(
			(p) =>
				p.name.toLowerCase().includes(searchQuery.toLowerCase()) ||
				p.category.toLowerCase().includes(searchQuery.toLowerCase())
		)
	);

	function resetForm() {
		formData = { name: '', category: '', price: '', stock: '' };
		editingId = null;
	}

	function openAdd() {
		resetForm();
		showAddModal = true;
	}

	function openEdit(product: (typeof products)[0]) {
		editingId = product.id;
		formData = {
			name: product.name,
			category: product.category,
			price: String(product.price),
			stock: String(product.stock)
		};
		showAddModal = true;
	}

	function handleSubmit(e: SubmitEvent) {
		e.preventDefault();
		if (!formData.name || !formData.category || !formData.price) return;

		if (editingId) {
			products = products.map((p) =>
				p.id === editingId
					? {
							...p,
							name: formData.name,
							category: formData.category,
							price: parseFloat(formData.price) || 0,
							stock: parseInt(formData.stock, 10) || 0
						}
					: p
			);
		} else {
			products = [
				...products,
				{
					id: String(Date.now()),
					name: formData.name,
					category: formData.category,
					price: parseFloat(formData.price) || 0,
					stock: parseInt(formData.stock, 10) || 0,
					status: 'active'
				}
			];
		}
		showAddModal = false;
		resetForm();
	}

	function deleteProduct(id: string) {
		products = products.filter((p) => p.id !== id);
	}

	function toggleStatus(id: string) {
		products = products.map((p) =>
			p.id === id ? { ...p, status: p.status === 'active' ? 'inactive' : 'active' } : p
		);
	}

	function stockLabel(stock: number): string {
		if (stock === 0) return 'Out of Stock';
		if (stock < 20) return 'Low Stock';
		return 'In Stock';
	}

	function stockColor(stock: number): string {
		if (stock === 0) return 'bg-red-100 text-red-800';
		if (stock < 20) return 'bg-yellow-100 text-yellow-800';
		return 'bg-green-100 text-green-800';
	}
</script>

<svelte:head>
	<title>Products - Coffee Admin</title>
</svelte:head>

<div class="space-y-6">
	<div class="flex items-center justify-between">
		<div>
			<h1 class="text-2xl font-bold text-foreground">Products</h1>
			<p class="mt-1 text-muted-foreground">{products.length} products total</p>
		</div>
		<Button onclick={openAdd}>
			<svg class="mr-2 h-4 w-4" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
				<path d="M12 5v14M5 12h14" />
			</svg>
			Add Product
		</Button>
	</div>

	<div class="flex items-center gap-4">
		<div class="relative flex-1 max-w-sm">
			<svg class="absolute left-3 top-1/2 h-4 w-4 -translate-y-1/2 text-muted-foreground" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
				<circle cx="11" cy="11" r="8" />
				<path d="m21 21-4.3-4.3" />
			</svg>
			<input
				type="text"
				bind:value={searchQuery}
				placeholder="Search products..."
				class="w-full rounded-md border border-border bg-background py-2 pl-10 pr-3 text-sm placeholder:text-muted-foreground focus:outline-none focus:ring-2 focus:ring-ring"
			/>
		</div>
	</div>

	<div class="rounded-lg border border-border bg-card">
		<div class="overflow-x-auto">
			<table class="w-full">
				<thead>
					<tr class="border-b border-border bg-muted/50">
						<th class="px-6 py-3 text-left text-xs font-medium uppercase tracking-wider text-muted-foreground">Name</th>
						<th class="px-6 py-3 text-left text-xs font-medium uppercase tracking-wider text-muted-foreground">Category</th>
						<th class="px-6 py-3 text-left text-xs font-medium uppercase tracking-wider text-muted-foreground">Price</th>
						<th class="px-6 py-3 text-left text-xs font-medium uppercase tracking-wider text-muted-foreground">Stock</th>
						<th class="px-6 py-3 text-left text-xs font-medium uppercase tracking-wider text-muted-foreground">Status</th>
						<th class="px-6 py-3 text-right text-xs font-medium uppercase tracking-wider text-muted-foreground">Actions</th>
					</tr>
				</thead>
				<tbody class="divide-y divide-border">
					{#each filteredProducts as product (product.id)}
						<tr class="hover:bg-muted/50">
							<td class="whitespace-nowrap px-6 py-4 text-sm font-medium text-foreground">{product.name}</td>
							<td class="whitespace-nowrap px-6 py-4 text-sm text-muted-foreground">{product.category}</td>
							<td class="whitespace-nowrap px-6 py-4 text-sm text-foreground">${product.price.toFixed(2)}</td>
							<td class="whitespace-nowrap px-6 py-4">
								<span class="inline-flex rounded-full px-2 py-1 text-xs font-medium {stockColor(product.stock)}">
									{product.stock} - {stockLabel(product.stock)}
								</span>
							</td>
							<td class="whitespace-nowrap px-6 py-4">
								<span class="inline-flex rounded-full px-2 py-1 text-xs font-medium {product.status === 'active' ? 'bg-green-100 text-green-800' : 'bg-muted text-muted-foreground'}">
									{product.status}
								</span>
							</td>
							<td class="whitespace-nowrap px-6 py-4 text-right">
								<div class="flex items-center justify-end gap-2">
									<button class="rounded-md p-1.5 text-muted-foreground transition-colors hover:bg-accent hover:text-accent-foreground" onclick={() => openEdit(product)} title="Edit">
										<svg class="h-4 w-4" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
											<path d="M17 3a2.85 2.83 0 1 1 4 4L7.5 20.5 2 22l1.5-5.5Z" />
											<path d="m15 5 4 4" />
										</svg>
									</button>
									<button class="rounded-md p-1.5 text-muted-foreground transition-colors hover:bg-accent hover:text-accent-foreground" onclick={() => toggleStatus(product.id)} title="Toggle Status">
										<svg class="h-4 w-4" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
											<path d="M12 2v20M2 12h20" />
										</svg>
									</button>
									<button class="rounded-md p-1.5 text-destructive transition-colors hover:bg-destructive/10" onclick={() => deleteProduct(product.id)} title="Delete">
										<svg class="h-4 w-4" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
											<path d="M3 6h18M19 6v14c0 1-1 2-2 2H7c-1 0-2-1-2-2V6M8 6V4c0-1 1-2 2-2h4c1 0 2 1 2 2v2" />
										</svg>
									</button>
								</div>
							</td>
						</tr>
					{:else}
						<tr>
							<td colspan="6" class="px-6 py-12 text-center text-muted-foreground">
								No products found
							</td>
						</tr>
					{/each}
				</tbody>
			</table>
		</div>
	</div>
</div>

{#if showAddModal}
	<!-- svelte-ignore a11y_no_static_element_interactions a11y_click_events_have_key_events -->
	<div
		class="fixed inset-0 z-50 flex items-center justify-center bg-black/50"
		onclick={() => { showAddModal = false; resetForm(); }}
	>
		<!-- svelte-ignore a11y_no_static_element_interactions a11y_click_events_have_key_events -->
		<div class="w-full max-w-md rounded-lg border border-border bg-card p-6" onclick={(e) => e.stopPropagation()}>
			<h2 class="text-lg font-semibold text-foreground mb-4">
				{editingId ? 'Edit Product' : 'Add New Product'}
			</h2>
			<form onsubmit={handleSubmit}>
				<div class="space-y-4">
					<div>
						<label class="mb-1 block text-sm font-medium text-foreground" for="product-name">Name</label>
						<input
							id="product-name"
							type="text"
							bind:value={formData.name}
							class="w-full rounded-md border border-border bg-background px-3 py-2 text-sm placeholder:text-muted-foreground focus:outline-none focus:ring-2 focus:ring-ring"
							placeholder="e.g., Americano"
							required
						/>
					</div>
					<div>
						<label class="mb-1 block text-sm font-medium text-foreground" for="product-category">Category</label>
						<input
							id="product-category"
							type="text"
							bind:value={formData.category}
							class="w-full rounded-md border border-border bg-background px-3 py-2 text-sm placeholder:text-muted-foreground focus:outline-none focus:ring-2 focus:ring-ring"
							placeholder="e.g., Coffee"
							required
						/>
					</div>
					<div class="grid grid-cols-2 gap-4">
						<div>
							<label class="mb-1 block text-sm font-medium text-foreground" for="product-price">Price ($)</label>
							<input
								id="product-price"
								type="number"
								step="0.01"
								min="0"
								bind:value={formData.price}
								class="w-full rounded-md border border-border bg-background px-3 py-2 text-sm placeholder:text-muted-foreground focus:outline-none focus:ring-2 focus:ring-ring"
								placeholder="4.00"
								required
							/>
						</div>
						<div>
							<label class="mb-1 block text-sm font-medium text-foreground" for="product-stock">Stock</label>
							<input
								id="product-stock"
								type="number"
								min="0"
								bind:value={formData.stock}
								class="w-full rounded-md border border-border bg-background px-3 py-2 text-sm placeholder:text-muted-foreground focus:outline-none focus:ring-2 focus:ring-ring"
								placeholder="100"
							/>
						</div>
					</div>
				</div>
				<div class="mt-6 flex gap-3">
					<button type="submit" class="flex-1 rounded-md bg-primary px-4 py-2 text-sm font-medium text-primary-foreground transition-colors hover:bg-primary/90">
						{editingId ? 'Save Changes' : 'Add Product'}
					</button>
					<button type="button" class="flex-1 rounded-md border border-border px-4 py-2 text-sm font-medium transition-colors hover:bg-muted" onclick={() => { showAddModal = false; resetForm(); }}>
						Cancel
					</button>
				</div>
			</form>
		</div>
	</div>
{/if}
