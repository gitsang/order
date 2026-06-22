<script lang="ts">
	import type { Snippet } from 'svelte';
	import { page } from '$app/stores';

	let { children }: { children: Snippet } = $props();

	const navItems = [
		{ href: '/admin', label: 'Dashboard', icon: 'home' },
		{ href: '/admin/products', label: 'Products', icon: 'package' },
		{ href: '/admin/orders', label: 'Orders', icon: 'shopping-bag' }
	];
</script>

<div class="flex h-screen bg-background">
	<!-- Sidebar -->
	<aside class="flex w-64 flex-col border-r border-border bg-card">
		<!-- Logo / Brand -->
		<div class="flex h-16 items-center border-b border-border px-6">
			<a href="/admin" class="flex items-center gap-2">
				<svg
					class="h-7 w-7 text-primary"
					viewBox="0 0 24 24"
					fill="none"
					stroke="currentColor"
					stroke-width="2"
					stroke-linecap="round"
					stroke-linejoin="round"
				>
					<path d="M18 8h1a4 4 0 010 8h-1" />
					<path d="M2 8h16v9a4 4 0 01-4 4H6a4 4 0 01-4-4V8z" />
					<path d="M6 1v3M10 1v3M14 1v3" />
				</svg>
				<span class="text-lg font-bold text-foreground">Coffee Admin</span>
			</a>
		</div>

		<!-- Navigation -->
		<nav class="flex-1 space-y-1 px-3 py-4">
			{#each navItems as item}
				{@const isActive =
					$page.url.pathname === item.href ||
					(item.href !== '/admin' && $page.url.pathname.startsWith(item.href))}
				<a
					href={item.href}
					class="flex items-center gap-3 rounded-md px-3 py-2 text-sm font-medium transition-colors {isActive
						? 'bg-primary text-primary-foreground'
						: 'text-muted-foreground hover:bg-accent hover:text-accent-foreground'}"
				>
					{#if item.icon === 'home'}
						<svg
							class="h-5 w-5"
							viewBox="0 0 24 24"
							fill="none"
							stroke="currentColor"
							stroke-width="2"
							stroke-linecap="round"
							stroke-linejoin="round"
						>
							<path
								d="M3 12l2-2m0 0l7-7 7 7M5 10v10a1 1 0 001 1h3m10-11l2 2m-2-2v10a1 1 0 01-1 1h-3m-6 0a1 1 0 001-1v-4a1 1 0 011-1h2a1 1 0 011 1v4a1 1 0 001 1m-6 0h6"
							/>
						</svg>
					{:else if item.icon === 'package'}
						<svg
							class="h-5 w-5"
							viewBox="0 0 24 24"
							fill="none"
							stroke="currentColor"
							stroke-width="2"
							stroke-linecap="round"
							stroke-linejoin="round"
						>
							<path
								d="M20 7l-8-4-8 4m16 0l-8 4m8-4v10l-8 4m0-10L4 7m8 4v10M4 7v10l8 4"
							/>
						</svg>
					{:else if item.icon === 'shopping-bag'}
						<svg
							class="h-5 w-5"
							viewBox="0 0 24 24"
							fill="none"
							stroke="currentColor"
							stroke-width="2"
							stroke-linecap="round"
							stroke-linejoin="round"
						>
							<path d="M16 11V7a4 4 0 00-8 0v4M5 9h14l1 12H4L5 9z" />
						</svg>
					{/if}
					{item.label}
				</a>
			{/each}
		</nav>

		<!-- Footer: Back to store -->
		<div class="border-t border-border p-4">
			<a
				href="/"
				class="flex items-center gap-3 rounded-md px-3 py-2 text-sm font-medium text-muted-foreground transition-colors hover:bg-accent hover:text-accent-foreground"
			>
				<svg
					class="h-5 w-5"
					viewBox="0 0 24 24"
					fill="none"
					stroke="currentColor"
					stroke-width="2"
					stroke-linecap="round"
					stroke-linejoin="round"
				>
					<path d="M13 7l5 5m0 0l-5 5m5-5H6" />
				</svg>
				Back to Store
			</a>
		</div>
	</aside>

	<!-- Main Content -->
	<main class="flex-1 overflow-y-auto">
		<div class="p-8">
			{@render children()}
		</div>
	</main>
</div>
