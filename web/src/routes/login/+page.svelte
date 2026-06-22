<script lang="ts">
	import { authApi } from '$lib/api';
	import { auth } from '$lib/stores';
	import { Button } from '$lib/components/ui/button';
	import { goto } from '$app/navigation';

	let username = $state('');
	let password = $state('');
	let isLoading = $state(false);
	let error = $state('');

	async function handleSubmit(e: SubmitEvent) {
		e.preventDefault();
		if (!username.trim() || !password.trim()) {
			error = 'Please enter username and password';
			return;
		}

		isLoading = true;
		error = '';

		try {
			const { token } = await authApi.login({ username: username.trim(), password });
			const user = await authApi.getMe();
			auth.login(token, user);

			if (user.role === 'admin') {
				goto('/admin');
			} else {
				goto('/');
			}
		} catch (err) {
			error = err instanceof Error ? err.message : 'Login failed. Please try again.';
		} finally {
			isLoading = false;
		}
	}
</script>

<div class="flex min-h-screen items-center justify-center bg-background px-4">
	<div class="w-full max-w-sm">
		<div class="mb-8 text-center">
			<svg
				class="mx-auto h-12 w-12 text-primary"
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
			<h1 class="mt-3 text-2xl font-bold text-foreground">Welcome Back</h1>
			<p class="mt-1 text-sm text-muted-foreground">Sign in to your account</p>
		</div>

		<div class="rounded-xl border border-border bg-card p-6 shadow-sm">
			<form onsubmit={handleSubmit} class="space-y-4">
				{#if error}
					<div class="rounded-lg bg-destructive/10 px-4 py-3 text-sm text-destructive">
						{error}
					</div>
				{/if}

				<div class="space-y-2">
					<label for="username" class="text-sm font-medium text-foreground">Username</label>
					<input
						id="username"
						type="text"
						bind:value={username}
						placeholder="Enter your username"
						autocomplete="username"
						disabled={isLoading}
						class="flex h-10 w-full rounded-md border border-input bg-background px-3 py-2 text-sm ring-offset-background placeholder:text-muted-foreground focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2 disabled:cursor-not-allowed disabled:opacity-50"
					/>
				</div>

				<div class="space-y-2">
					<label for="password" class="text-sm font-medium text-foreground">Password</label>
					<input
						id="password"
						type="password"
						bind:value={password}
						placeholder="Enter your password"
						autocomplete="current-password"
						disabled={isLoading}
						class="flex h-10 w-full rounded-md border border-input bg-background px-3 py-2 text-sm ring-offset-background placeholder:text-muted-foreground focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2 disabled:cursor-not-allowed disabled:opacity-50"
					/>
				</div>

				<Button type="submit" class="w-full" size="lg" disabled={isLoading}>
					{#if isLoading}
						Signing in...
					{:else}
						Sign In
					{/if}
				</Button>
			</form>

			<div class="mt-4 text-center text-sm text-muted-foreground">
				Don't have an account?
				<a href="/register" class="font-medium text-primary hover:underline">Sign up</a>
			</div>
		</div>
	</div>
</div>
