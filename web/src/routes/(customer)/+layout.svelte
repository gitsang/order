<script lang="ts">
	import type { Snippet } from 'svelte';
	import { page } from '$app/stores';
	import { Home, ShoppingCart, ClipboardList } from 'lucide-svelte';
	import { cart } from '$lib/stores';
	import { cn } from '$lib/utils';

	let { children }: { children: Snippet } = $props();

	const tabs = [
		{ href: '/', label: 'Home', icon: Home },
		{ href: '/cart', label: 'Cart', icon: ShoppingCart },
		{ href: '/order', label: 'Orders', icon: ClipboardList }
	];
</script>

<div class="flex min-h-screen flex-col bg-background">
	<main class="flex-1 pb-20">
		{@render children()}
	</main>

	<nav
		class="fixed bottom-0 left-0 right-0 z-50 border-t border-border bg-card/95 backdrop-blur supports-[backdrop-filter]:bg-card/60"
	>
		<div class="mx-auto flex h-16 max-w-lg items-center justify-around px-4">
			{#each tabs as tab (tab.href)}
				{@const isActive =
					tab.href === '/'
						? $page.url.pathname === '/'
						: $page.url.pathname.startsWith(tab.href)}
				<a
					href={tab.href}
					class={cn(
						'flex flex-col items-center justify-center gap-1 rounded-lg px-3 py-2 text-xs font-medium transition-colors',
						isActive
							? 'text-primary'
							: 'text-muted-foreground hover:text-foreground'
					)}
				>
					<div class="relative">
						<tab.icon class="h-5 w-5" strokeWidth={isActive ? 2.5 : 2} />
						{#if tab.href === '/cart' && $cart.totalItems > 0}
							<span
								class="absolute -right-2 -top-2 flex h-4 w-4 items-center justify-center rounded-full bg-primary text-[10px] font-bold text-primary-foreground"
							>
								{$cart.totalItems > 99 ? '99+' : $cart.totalItems}
							</span>
						{/if}
					</div>
					<span>{tab.label}</span>
				</a>
			{/each}
		</div>
	</nav>
</div>
