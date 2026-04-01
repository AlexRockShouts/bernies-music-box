<script>
  import { goto } from '$app/navigation';
  import { onMount } from 'svelte';
  import { createTask, getHistory } from '$lib/api.js';
  let history = [];
  let prompt = '';
  let style = 'indie rock, electric guitar, driving drums, male vocal, E minor';
  let duration = 30;
  let loading = false;
  let error = '';
  onMount(async () => {
    const token = localStorage.getItem('token');
    if (!token) {
      goto('/login');
      return;
    }
    await loadHistory();
  });
  async function loadHistory() {
    try {
      history = await getHistory();
    } catch (e) {
      error = e.message;
    }
  }
  async function submitPrompt() {
    if (!prompt.trim()) return;
    loading = true;
    error = '';
    try {
      const task = await createTask({ prompt, style, duration });
      console.log('Created task:', task);
      prompt = '';
      await loadHistory();
    } catch (e) {
      error = e.message;
    } finally {
      loading = false;
    }
  }
</script>

{#if error}
  <div class="error">{error}</div>
{/if}

<h1>Bernie's Music Box</h1>
<button onclick={() => { localStorage.removeItem('token'); goto('/login'); }} >Logout</button>

<h2>Historical Prompts</h2>
{#if history.length === 0}
  <p>No tasks yet.</p>
{:else}
  <div class="history">
    {#each history as task}
      <div class="task">
        <h3>ID: {task.id}</h3>
        <p><strong>Prompt:</strong> {task.prompt}</p>
        <p><strong>Status:</strong> {task.status}</p>
        <p><strong>Created:</strong> {new Date(task.createdAt).toLocaleString()}</p>
      </div>
    {/each}
  </div>
{/if}

<h2>New Prompt</h2>
<form on:submit|preventDefault={submitPrompt}>
  <textarea bind:value={prompt} placeholder="Enter your lyrics/prompt..." rows="4"></textarea>
  <div class="params">
    <label>
      Style:
      <select bind:value={style}>
        <option value="indie rock, electric guitar, driving drums, male vocal, E minor">Indie Rock</option>
        <option value="pop, synths, upbeat, female vocal, C major">Pop</option>
        <option value="jazz, piano, smooth, sax, A minor">Jazz</option>
      </select>
    </label>
    <label>
      Duration (s):
      <input type="number" bind:value={duration} min="10" max="120" />
    </label>
  </div>
  <button type="submit" disabled={loading}>{loading ? 'Generating...' : 'Generate Music'}</button>
</form>

<style>
  :global(body) {
    font-family: Arial, sans-serif;
    max-width: 1200px;
    margin: 0 auto;
    padding: 1rem;
  }
  h1, h2 {
    color: #333;
  }
  .error {
    color: red;
    background: #fee;
    padding: 1rem;
    border-radius: 4px;
  }
  textarea {
    width: 100%;
    padding: 0.5rem;
    border: 1px solid #ccc;
    border-radius: 4px;
    font-size: 1rem;
  }
  .params {
    display: flex;
    gap: 1rem;
    margin: 1rem 0;
  }
  label {
    flex: 1;
    display: flex;
    flex-direction: column;
  }
  select, input {
    padding: 0.5rem;
    border: 1px solid #ccc;
    border-radius: 4px;
  }
  button {
    padding: 0.75rem 1.5rem;
    background: #007acc;
    color: white;
    border: none;
    border-radius: 4px;
    cursor: pointer;
  }
  button:disabled {
    background: #ccc;
  }
  .history {
    display: grid;
    gap: 1rem;
    margin-bottom: 2rem;
  }
  .task {
    border: 1px solid #ddd;
    padding: 1rem;
    border-radius: 8px;
  }
</style>