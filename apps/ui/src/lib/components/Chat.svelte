<script lang="ts">
	export let messages: Array<{ id: number; text: string; sender: "user" | "assistant" }> = [];

	let inputText = "";
	let messageContainer: HTMLDivElement;

	const sendMessage = () => {
		if (inputText.trim() === "") return;

		messages = [
			...messages,
			{ id: messages.length, text: inputText, sender: "user" }
		];
		inputText = "";

		setTimeout(() => {
			if (messageContainer) {
				messageContainer.scrollTop = messageContainer.scrollHeight;
			}
		}, 0);
	};

	const handleKeyDown = (e: KeyboardEvent) => {
		if (e.key === "Enter" && !e.shiftKey) {
			e.preventDefault();
			sendMessage();
		}
	};
</script>

<div class="flex flex-col h-full">
	<!-- Messages Container -->
	<div bind:this={messageContainer} class="flex-1 overflow-y-auto p-3 md:p-4 space-y-4">
		{#each messages as message (message.id)}
			<div class="flex {message.sender === 'user' ? 'justify-end' : 'justify-start'}">
				<div
					class="max-w-xs px-3 md:px-4 py-2 rounded-lg text-sm md:text-base {message.sender === 'user'
						? 'bg-[#77B6EA] text-white rounded-br-none'
						: 'bg-gray-200 text-gray-800 rounded-bl-none'}"
				>
					{message.text}
				</div>
			</div>
		{/each}
	</div>

	<!-- Input Area -->
	<div class="border-t border-gray-300 p-3 md:p-4">
		<div class="flex gap-2 flex-col sm:flex-row">
			<textarea
				bind:value={inputText}
				on:keydown={handleKeyDown}
				class="flex-1 p-2 md:p-3 border border-gray-300 rounded-lg resize-none focus:outline-none focus:ring-2 focus:ring-[#77B6EA] text-sm md:text-base max-h-24"
				placeholder="Type your message... (Shift+Enter for new line)"
				rows="1"
			/>
			<button
				on:click={sendMessage}
				class="px-3 md:px-4 py-2 md:py-3 bg-[#77B6EA] text-white rounded-lg hover:bg-[#5a9ecb] transition-colors duration-200 font-semibold text-sm md:text-base whitespace-nowrap"
			>
				Send
			</button>
		</div>
	</div>
</div>
