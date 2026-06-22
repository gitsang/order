<script lang="ts">
	import { onMount } from 'svelte';
	import { Button } from '$lib/components/ui/button';
	import { categoriesApi } from '$lib/api';
	import type { Category } from '$lib/api/types';

	let categories = $state<Category[]>([]);
	let loading = $state(true);
	let error = $state<string | null>(null);

	let showAddModal = $state(false);
	let editingId = $state<string | null>(null);
	let searchQuery = $state('');
	let submitting = $state(false);
	let deleteConfirmId = $state<string | null>(null);

	let formData = $state({ name: '', sort_order: '' });

	onMount(async () => {
		await fetchCategories();
	});

	async function fetchCategories() {
		loading = true;
		error = null;
		try {
			categories = await categoriesApi.list();
		} catch (e) {
			error = e instanceof Error ? e.message : 'Failed to load categories';
		} finally {
			loading = false;
		}
	}

	let filteredCategories = $derived(
		categories.filter((c) => c.name.toLowerCase().includes(searchQuery.toLowerCase()))
	);

	function resetForm() {
		formData = { name: '', sort_order: '' };
		editingId = null;
	}

	function openAdd() {
		resetForm();
		showAddModal = true;
	}

	function openEdit(category: Category) {
		editingId = category.id;
		formData = {
			name: category.name,
			sort_order: String(category.sort_order)
		};
		showAddModal = true;
	}

	async function handleSubmit(e: SubmitEvent) {
		e.preventDefault();
		if (!formData.name) return;

		submitting = true;
		try {
			if (editingId) {
				const updated = await categoriesApi.update(editingId, {
					name: formData.name,
					sort_order: parseInt(formData.sort_order, 10) || 0
				});
				categories = categories.map((c) => (c.id === editingId ? updated : c));
			} else {
				const created = await categoriesApi.create({
					name: formData.name,
					sort_order: parseInt(formData.sort_order, 10) || 0
				});
				categories = [...categories, created];
			}
			showAddModal = false;
			resetForm();
		} catch (e) {
			error = e instanceof Error ? e.message : 'Failed to save category';
		} finally {
			submitting = false;
		}
	}

	async function deleteCategory(id: string) {
		try {
			await categoriesApi.delete(id);
			categories = categories.filter((c) => c.id !== id);
			deleteConfirmId = null;
		} catch (e) {
			error = e instanceof Error ? e.message : 'Failed to delete category';
		}
	}
</script>

<svelte:head>
	<title>Categories - Coffee Admin</title>
</svelte:head>

<div class="space-y-6">
	<div class="flex items-center justify-between">
		<div>
			<h1 class="text-2xl font-bold text-foreground">Categories</h1>
			<p class="mt-1 text-muted-foreground">{loading ? 'Loading...' : `${categories.length} categories total`}</p>
		</div>
		<Button onclick={openAdd}>
			<svg class="mr-2 h-4 w-4" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
				<path d="M12 5v14M5 12h14" />
			</svg>
			Add Category
		</Button>
	</div>

	{#if error}
		<div class="rounded-md bg-destructive/10 p-4 text-sm text-destructive">
			{error}
			<button class="ml-2 underline" onclick={() => { error = null; fetchCategories(); }}>Retry</button>
		</div>
	{/if}

	<div class="flex items-center gap-4">
		<div class="relative flex-1 max-w-sm">
			<svg class="absolute left-3 top-1/2 h-4 w-4 -translate-y-1/2 text-muted-foreground" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
				<circle cx="11" cy="11" r="8" />
				<path d="m21 21-4.3-4.3" />
			</svg>
			<input
				type="text"
				bind:value={searchQuery}
				placeholder="Search categories..."
				class="w-full rounded-md border border-border bg-background py-2 pl-10 pr-3 text-sm placeholder:text-muted-foreground focus:outline-none focus:ring-2 focus:ring-ring"
			/>
		</div>
	</div>

	<div class="rounded-lg border border-border bg-card">
		<div class="overflow-x-auto">
			{#if loading}
				<div class="flex items-center justify-center py-12">
					<svg class="h-8 w-8 animate-spin text-primary" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
						<path d="M12 2v4m0 12v4M4.93 4.93l2.83 2.83m8.48 8.48l2.83 2.83M2 12h4m12 0h4M4.93 19.07l2.83-2.83m8.48-8.48l2.83-2.83" />
					</svg>
				</div>
			{:else}
				<table class="w-full">
					<thead>
						<tr class="border-b border-border bg-muted/50">
							<th class="px-6 py-3 text-left text-xs font-medium uppercase tracking-wider text-muted-foreground">Name</th>
							<th class="px-6 py-3 text-left text-xs font-medium uppercase tracking-wider text-muted-foreground">Sort Order</th>
							<th class="px-6 py-3 text-right text-xs font-medium uppercase tracking-wider text-muted-foreground">Actions</th>
						</tr>
					</thead>
					<tbody class="divide-y divide-border">
						{#each filteredCategories as category (category.id)}
							<tr class="hover:bg-muted/50">
								<td class="whitespace-nowrap px-6 py-4 text-sm font-medium text-foreground">{category.name}</td>
								<td class="whitespace-nowrap px-6 py-4 text-sm text-muted-foreground">{category.sort_order}</td>
								<td class="whitespace-nowrap px-6 py-4 text-right">
									<div class="flex items-center justify-end gap-2">
										<button class="rounded-md p-1.5 text-muted-foreground transition-colors hover:bg-accent hover:text-accent-foreground" onclick={() => openEdit(category)} title="Edit">
											<svg class="h-4 w-4" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
												<path d="M17 3a2.85 2.83 0 1 1 4 4L7.5 20.5 2 22l1.5-5.5Z" />
												<path d="m15 5 4 4" />
											</svg>
										</button>
										<button class="rounded-md p-1.5 text-destructive transition-colors hover:bg-destructive/10" onclick={() => (deleteConfirmId = category.id)} title="Delete">
											<svg class="h-4 w-4" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
												<path d="M3 6h18M19 6v14c0 1-1 2-2 2H7c-1 0-2-1-2-2V6M8 6V4c0-1 1-2 2-2h4c1 0 2 1 2 2v2" />
											</svg>
										</button>
									</div>
								</td>
							</tr>
						{:else}
							<tr>
								<td colspan="3" class="px-6 py-12 text-center text-muted-foreground">
									{searchQuery ? 'No categories match your search' : 'No categories found'}
								</td>
							</tr>
						{/each}
					</tbody>
				</table>
			{/if}
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
				{editingId ? 'Edit Category' : 'Add New Category'}
			</h2>
			<form onsubmit={handleSubmit}>
				<div class="space-y-4">
					<div>
						<label class="mb-1 block text-sm font-medium text-foreground" for="category-name">Name</label>
						<input
							id="category-name"
							type="text"
							bind:value={formData.name}
							class="w-full rounded-md border border-border bg-background px-3 py-2 text-sm placeholder:text-muted-foreground focus:outline-none focus:ring-2 focus:ring-ring"
							placeholder="e.g., Coffee"
							required
						/>
					</div>
					<div>
						<label class="mb-1 block text-sm font-medium text-foreground" for="category-sort">Sort Order</label>
						<input
							id="category-sort"
							type="number"
							min="0"
							bind:value={formData.sort_order}
							class="w-full rounded-md border border-border bg-background px-3 py-2 text-sm placeholder:text-muted-foreground focus:outline-none focus:ring-2 focus:ring-ring"
							placeholder="0"
						/>
					</div>
				</div>
				<div class="mt-6 flex gap-3">
					<button
						type="submit"
						disabled={submitting}
						class="flex-1 rounded-md bg-primary px-4 py-2 text-sm font-medium text-primary-foreground transition-colors hover:bg-primary/90 disabled:opacity-50"
					>
						{submitting ? 'Saving...' : editingId ? 'Save Changes' : 'Add Category'}
					</button>
					<button type="button" class="flex-1 rounded-md border border-border px-4 py-2 text-sm font-medium transition-colors hover:bg-muted" onclick={() => { showAddModal = false; resetForm(); }}>
						Cancel
					</button>
				</div>
			</form>
		</div>
	</div>
{/if}

{#if deleteConfirmId}
	<!-- svelte-ignore a11y_no_static_element_interactions a11y_click_events_have_key_events -->
	<div
		class="fixed inset-0 z-50 flex items-center justify-center bg-black/50"
		onclick={() => (deleteConfirmId = null)}
	>
		<!-- svelte-ignore a11y_no_static_element_interactions a11y_click_events_have_key_events -->
		<div class="w-full max-w-sm rounded-lg border border-border bg-card p-6" onclick={(e) => e.stopPropagation()}>
			<h2 class="text-lg font-semibold text-foreground mb-2">Delete Category</h2>
			<p class="text-sm text-muted-foreground mb-4">Are you sure you want to delete this category? Products in this category may be affected.</p>
			<div class="flex gap-3">
				<button
					class="flex-1 rounded-md bg-destructive px-4 py-2 text-sm font-medium text-destructive-foreground transition-colors hover:bg-destructive/90"
					onclick={() => deleteCategory(deleteConfirmId!)}
				>
					Delete
				</button>
				<button
					class="flex-1 rounded-md border border-border px-4 py-2 text-sm font-medium transition-colors hover:bg-muted"
					onclick={() => (deleteConfirmId = null)}
				>
					Cancel
				</button>
			</div>
		</div>
	</div>
{/if}
