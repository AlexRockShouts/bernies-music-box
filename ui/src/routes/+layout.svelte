<script>
	import "../app.css";
	let { children } = $props();
</script>

<div class="app">
	<header>
		<h1>Berne's Music Box</h1>
	</header>

	<main>
		{@render children()}
	</main>

	<footer>
		<p>&copy; 2026 Berne's Music Box</p>
	</footer>
</div>

<style>
	:global(body) {
		margin: 0;
		font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, Oxygen, Ubuntu, Cantarell, 'Open Sans', 'Helvetica Neue', sans-serif;
		background-color: #0f172a;
		color: #f8fafc;
	}

	.app {
		display: flex;
		flex-direction: column;
		min-height: 100vh;
	}

	header {
		padding: 1rem 2rem;
		background-color: #1e293b;
		border-bottom: 1px solid #334155;
	}

	h1 {
		margin: 0;
		font-size: 1.5rem;
		color: #38bdf8;
	}

	main {
		flex: 1;
		padding: 2rem;
		max-width: 1200px;
		margin: 0 auto;
		width: 100%;
		box-sizing: border-box;
	}

	footer {
		padding: 1rem;
		text-align: center;
		font-size: 0.875rem;
		color: #94a3b8;
		border-top: 1px solid #334155;
	}
</style>
