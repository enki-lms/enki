<script lang="ts">
  import { api } from "$lib/api";
  import { marked } from "marked";

  let { messages = $bindable([]), additionalContext = "" } = $props<{
    messages: Array<{ id: number; text: string; sender: "user" | "assistant" }>;
    additionalContext?: string;
  }>();

  let inputText = $state("");
  let messageContainer: HTMLDivElement;
  let loading = $state(false);

  const sendMessage = async () => {
    if (inputText.trim() === "" || loading) return;

    const userMessageText = inputText;
    inputText = "";

    const newUserMessage = {
      id: messages.length,
      text: userMessageText,
      sender: "user" as const,
    };
    messages = [...messages, newUserMessage];

    loading = true;

    // Scroll to bottom
    setTimeout(() => {
      if (messageContainer) {
        messageContainer.scrollTop = messageContainer.scrollHeight;
      }
    }, 0);

    try {
      const apiMessages = messages.map((m) => ({
        role: m.sender,
        content: m.text,
      }));

      if (additionalContext) {
        const lastMsg = apiMessages[apiMessages.length - 1];
        lastMsg.content = `[Context: ${additionalContext}]\n\nUser Question: ${lastMsg.content}`;
      }

      const response = await api.chat(apiMessages);

      messages = [
        ...messages,
        { id: messages.length, text: response.response, sender: "assistant" },
      ];
    } catch (e) {
      messages = [
        ...messages,
        {
          id: messages.length,
          text: "Error: Failed to get response from AI.",
          sender: "assistant",
        },
      ];
      console.error(e);
    } finally {
      loading = false;
      setTimeout(() => {
        if (messageContainer) {
          messageContainer.scrollTop = messageContainer.scrollHeight;
        }
      }, 0);
    }
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
  <div
    bind:this={messageContainer}
    class="flex-1 overflow-y-auto p-3 md:p-4 space-y-4"
  >
    {#each messages as message (message.id)}
      <div
        class="flex {message.sender === 'user'
          ? 'justify-end'
          : 'justify-start'}"
      >
        <div
          class="max-w-2xl px-3 md:px-4 py-2 rounded-lg text-sm md:text-base {message.sender ===
          'user'
            ? 'bg-[#77B6EA] text-white rounded-br-none'
            : 'bg-gray-200 text-gray-800 rounded-bl-none'}"
        >
          {#if message.sender === "assistant"}
            <div
              class="prose prose-sm max-w-none prose-headings:mb-2 prose-p:mb-2 prose-pre:bg-gray-800 prose-pre:text-white"
            >
              {@html marked.parse(message.text)}
            </div>
          {:else}
            {message.text}
          {/if}
        </div>
      </div>
    {/each}
  </div>

  <!-- Input Area -->
  <div class="border-t border-gray-300 p-3 md:p-4">
    <div class="flex gap-2 flex-col sm:flex-row">
      <textarea
        bind:value={inputText}
        onkeydown={handleKeyDown}
        class="flex-1 p-2 md:p-3 border border-gray-300 rounded-lg resize-none focus:outline-none focus:ring-2 focus:ring-[#77B6EA] text-sm md:text-base max-h-24 disabled:bg-gray-100 disabled:text-gray-500"
        placeholder={loading
          ? "Thinking..."
          : "Type your message... (Shift+Enter for new line)"}
        rows="1"
        disabled={loading}
      />
      <button
        onclick={sendMessage}
        disabled={loading}
        class="px-3 md:px-4 py-2 md:py-3 bg-[#77B6EA] text-white rounded-lg hover:bg-[#5a9ecb] transition-colors duration-200 font-semibold text-sm md:text-base whitespace-nowrap disabled:bg-gray-400 disabled:cursor-not-allowed"
      >
        {loading ? "..." : "Send"}
      </button>
    </div>
  </div>
</div>
