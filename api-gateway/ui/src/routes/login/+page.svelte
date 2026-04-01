<script>
  import { goto } from '$app/navigation';
  let username = '';
  let password = '';
  let error = '';
  async function login() {
    if (username === 'demo' && password === 'demo') {
      localStorage.setItem('token', 'fake-jwt-token');
      goto('/');
    } else {
      error = 'Invalid credentials. Use demo/demo';
    }
  }
</script>

<div class="login">
  <h1>Login to Bernie's Music Box</h1>
  {#if error}
    <p class="error">{error}</p>
  {/if}
  <form on:submit|preventDefault={login}>
    <label>
      Username:
      <input bind:value={username} required />
    </label>
    <label>
      Password:
      <input bind:value={password} type="password" required />
    </label>
    <button type="submit">Login</button>
  </form>
</div>

<style>
  .login {
    max-width: 400px;
    margin: 100px auto;
    padding: 2rem;
    border: 1px solid #ccc;
    border-radius: 8px;
  }
  label {
    display: block;
    margin-bottom: 1rem;
  }
  input {
    width: 100%;
    padding: 0.5rem;
    margin-top: 0.25rem;
  }
  button {
    width: 100%;
    padding: 0.75rem;
    background: #007acc;
    color: white;
    border: none;
    border-radius: 4px;
  }
  .error {
    color: red;
  }
</style>