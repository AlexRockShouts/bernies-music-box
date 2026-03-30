<script lang="ts">
	import { onMount } from 'svelte';

	interface Task {
		id: string;
		prompt: string;
		status: string;
		result_url?: string;
		created_at: string;
	}

	let prompt = $state("");
	let tasks = $state<Task[]>([]);
	let isGenerating = $state(false);

	const API_BASE = "http://localhost:8080";

	async function fetchHistory() {
		try {
			const res = await fetch(`${API_BASE}/history`);
			if (res.ok) {
				tasks = await res.json();
			}
		} catch (e) {
			console.error("Failed to fetch history", e);
		}
	}

	async function generate() {
		if (!prompt || isGenerating) return;
		isGenerating = true;
		try {
			const res = await fetch(`${API_BASE}/tasks`, {
				method: 'POST',
				headers: { 'Content-Type': 'application/json' },
				body: JSON.stringify({ prompt })
			});
			if (res.ok) {
				const newTask = await res.json();
				tasks = [newTask, ...tasks];
				prompt = "";
				// Start polling for this task
				pollTask(newTask.id);
			}
		} catch (e) {
			console.error("Failed to generate", e);
		} finally {
			isGenerating = false;
		}
	}

	async function pollTask(id: string) {
		const interval = setInterval(async () => {
			try {
				const res = await fetch(`${API_BASE}/tasks/${id}`);
				if (res.ok) {
					const updatedTask = await res.json();
					const index = tasks.findIndex(t => t.id === id);
					if (index !== -1) {
						tasks[index] = updatedTask;
					}
					if (updatedTask.status === 'completed' || updatedTask.status === 'failed') {
						clearInterval(interval);
					}
				}
			} catch (e) {
				console.error("Polling error", e);
				clearInterval(interval);
			}
		}, 3000);
	}

	onMount(fetchHistory);
</script>

<div class="container">
	<section class="prompt-editor">
		<h2>Prompt Editor</h2>
		<textarea 
			bind:value={prompt} 
			placeholder="Describe the music you want to generate..."
			rows="4"
		></textarea>
		<button onclick={generate} disabled={isGenerating || !prompt}>
			{isGenerating ? "Generating..." : "Generate Music"}
		</button>
	</section>

	<section class="history">
		<h2>History & Results</h2>
		<div class="task-list">
			{#each tasks as task (task.id)}
				<div class="task-card">
					<div class="task-info">
						<p class="prompt-text">"{task.prompt}"</p>
						<span class="status-badge {task.status}">{task.status}</span>
						<small>{new Date(task.created_at).toLocaleString()}</small>
					</div>
					
					{#if task.status === 'completed' && task.result_url}
						<div class="result-actions">
							<div class="player">
								<audio controls src={task.result_url}>
									Your browser does not support the audio element.
								</audio>
							</div>
							<div class="midi-preview">
								<span class="icon">🎹</span> MIDI Preview
							</div>
							<a href={task.result_url} download class="download-btn">
								Download
							</a>
						</div>
					{/if}
				</div>
			{/each}
		</div>
	</section>
</div>

<style>
	.container {
		display: grid;
		grid-template-columns: 1fr;
		gap: 2rem;
	}

	section {
		background-color: #1e293b;
		padding: 1.5rem;
		border-radius: 0.5rem;
		border: 1px solid #334155;
	}

	h2 {
		margin-top: 0;
		font-size: 1.25rem;
		margin-bottom: 1rem;
		color: #94a3b8;
	}

	textarea {
		width: 100%;
		background-color: #0f172a;
		color: white;
		border: 1px solid #334155;
		border-radius: 0.25rem;
		padding: 0.75rem;
		font-size: 1rem;
		margin-bottom: 1rem;
		resize: vertical;
		box-sizing: border-box;
	}

	button {
		background-color: #0284c7;
		color: white;
		border: none;
		padding: 0.75rem 1.5rem;
		border-radius: 0.25rem;
		font-weight: 600;
		cursor: pointer;
		transition: background-color 0.2s;
	}

	button:hover:not(:disabled) {
		background-color: #0369a1;
	}

	button:disabled {
		opacity: 0.5;
		cursor: not-allowed;
	}

	.task-list {
		display: flex;
		flex-direction: column;
		gap: 1rem;
	}

	.task-card {
		background-color: #334155;
		padding: 1rem;
		border-radius: 0.4rem;
	}

	.task-info {
		display: flex;
		align-items: center;
		gap: 1rem;
		margin-bottom: 1rem;
		flex-wrap: wrap;
	}

	.prompt-text {
		margin: 0;
		font-weight: 500;
		flex: 1;
		min-width: 200px;
	}

	.status-badge {
		font-size: 0.75rem;
		padding: 0.25rem 0.5rem;
		border-radius: 1rem;
		text-transform: uppercase;
		font-weight: 700;
	}

	.status-badge.pending { background-color: #64748b; }
	.status-badge.processing { background-color: #eab308; color: black; }
	.status-badge.completed { background-color: #22c55e; }
	.status-badge.failed { background-color: #ef4444; }

	.result-actions {
		display: flex;
		align-items: center;
		gap: 1.5rem;
		background-color: #1e293b;
		padding: 0.75rem;
		border-radius: 0.25rem;
		flex-wrap: wrap;
	}

	.player { flex: 1; }
	audio { width: 100%; height: 32px; }

	.midi-preview {
		font-size: 0.875rem;
		color: #cbd5e1;
		display: flex;
		align-items: center;
		gap: 0.4rem;
	}

	.download-btn {
		color: #38bdf8;
		text-decoration: none;
		font-size: 0.875rem;
		font-weight: 600;
	}

	.download-btn:hover { text-decoration: underline; }
</style>
