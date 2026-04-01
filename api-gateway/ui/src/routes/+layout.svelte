<script>
  import { page } from '$app/stores';
  import { goto } from '$app/navigation';
  $: if ($page.url.pathname !== '/login' && !localStorage.getItem('token')) {
    goto('/login');
  }
</script>

<header>
  <nav>
    {#if $page.url.pathname === '/login'}
      <h1>Bernie's Music Box</h1>
    {:else}
      <a href="/">Dashboard</a>
      <button onclick={() => { localStorage.removeItem('token'); goto('/login'); }}>Logout</button>
    {/if}
  </nav>
</header>

<main>
  <slot />
</main>

<style>
  header {
    background: #007acc;
    color: white;
    padding: 1rem;
  }
  nav {
    max-width: 1200px;
    margin: 0 auto;
    display: flex;
    justify-content: space-between;
    align-items: center;
  }
  a {
    color: white;
    text-decoration: none;
    font-size: 1.2rem;
  }
  button {
    background: none;
    color: white;
    border: 1px solid white;
    padding: 0.5rem 1rem;
    border-radius: 4px;
    cursor: pointer;
  }
  main {
    max-width: 1200px;
    margin: 0 auto;
    padding: 1rem;
  }
</style>