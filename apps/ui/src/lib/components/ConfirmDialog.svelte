<script lang="ts">
  import { createEventDispatcher } from "svelte";
  import { fade, scale } from "svelte/transition";

  export let isOpen = false;
  export let title = "Confirm Delete";
  export let message =
    "Are you sure you want to delete this item? This action cannot be undone.";
  export let confirmText = "Delete";
  export let cancelText = "Cancel";
  export let isDangerous = true;
  export let isLoading = false;

  const dispatch = createEventDispatcher<{
    confirm: void;
    cancel: void;
  }>();

  const handleConfirm = () => {
    dispatch("confirm");
  };

  const handleCancel = () => {
    dispatch("cancel");
    isOpen = false;
  };

  const handleKeydown = (e: KeyboardEvent) => {
    if (e.key === "Escape") {
      handleCancel();
    }
  };
</script>

{#if isOpen}
  <div
    class="fixed inset-0 bg-black/30 flex items-center justify-center z-50 p-4"
    on:click={handleCancel}
    on:keydown={handleKeydown}
    role="button"
    tabindex="0"
    transition:fade={{ duration: 200 }}
  >
    <div
      class="bg-white rounded-2xl shadow-xl max-w-md w-full p-6"
      on:click|stopPropagation
      role="alertdialog"
      aria-modal="true"
      transition:scale={{ duration: 200, start: 0.95 }}
    >
      <h2 class="text-xl font-semibold text-gray-900 mb-2">{title}</h2>
      <p class="text-gray-600 mb-6">{message}</p>

      <div class="flex gap-3 justify-end">
        <button
          type="button"
          on:click={handleCancel}
          class="px-5 py-2 rounded-full border-2 border-gray-300 text-gray-700 font-medium hover:bg-gray-50 transition-colors"
          disabled={isLoading}
        >
          {cancelText}
        </button>
        <button
          type="button"
          on:click={handleConfirm}
          class="px-5 py-2 rounded-full font-medium text-white transition-colors
						{isDangerous ? 'bg-red-500 hover:bg-red-600' : 'bg-sky-500 hover:bg-sky-600'}
						disabled:opacity-50"
          disabled={isLoading}
        >
          {#if isLoading}
            <span
              class="inline-block w-4 h-4 border-2 border-white border-t-transparent rounded-full animate-spin mr-2"
            ></span>
          {/if}
          {confirmText}
        </button>
      </div>
    </div>
  </div>
{/if}
