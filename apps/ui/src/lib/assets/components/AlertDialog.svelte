<script lang="ts">
    // Props are now driven EXTERNALLY (by the store via the layout)
    export let isOpen: boolean; // Controlled by store
    export let title: string;
    export let description: string;
    export let confirmText: string;
    export let cancelText: string | null; // Null means no cancel button
    export let showAction: boolean;
    export let isDangerous: boolean;
    export let isLoading: boolean;

    // Callbacks to interact with the store
    export let doConfirm: () => void; // Function to call when confirm is clicked
    export let doCancel: () => void; // Function to call when cancel/close happens

    let dialog: HTMLDialogElement;

    // Reactive statement to control the native dialog based on the isOpen prop
    $: if (dialog && isOpen !== dialog.open) {
        if (isOpen) {
            dialog.showModal();
        } else {
            // Ensure we only call close if it's actually open, prevent errors
            if (dialog.open) {
                dialog.close();
            }
        }
    }

    // When the dialog's native 'close' event fires (ESC, backdrop click, or dialog.close())
    // We need to inform the controlling logic (the store) via doCancel.
    // Check if it's already being closed programmatically by the store setting isOpen=false
    // If isOpen is still true when 'close' fires, it means it was an external close (ESC/backdrop).
    function handleDialogCloseEvent() {
        if (isOpen) {
            doCancel(); // Trigger the cancel logic in the store
        }
        // If isOpen is already false, the store initiated the close (via confirm/cancel button),
        // so we don't need to call doCancel again.
    }

    // Handle backdrop clicks specifically to trigger cancel
    function handleBackdropClick(event: MouseEvent) {
        if (event.target === dialog) {
            doCancel();
        }
    }

</script>

<dialog
        bind:this={dialog}
        on:close={handleDialogCloseEvent}
        on:click={handleBackdropClick}
        class="backdrop:bg-black/50 backdrop:backdrop-blur-xs p-0 bg-opacity-0 bg-[transparent] border-0 outline-hidden {$$props.class ?? ''}"
        aria-labelledby="dialog-title"
        aria-describedby="dialog-description"
        role={cancelText ? "alertdialog" : "alert"}
>
    {#if isOpen}
        <div class="bg-blue-50 rounded-lg max-w-md w-full p-6 relative">
            <h2 id="dialog-title" class="text-xl font-semibold text-slate-900 mb-2">{title}</h2>
            <div class="flex items-center gap-3 mb-6">
                <p id="dialog-description" class="text-slate-700">{description}</p>
                {#if isLoading}
                    <div class="loader"></div>
                {/if}
            </div>

            {#if showAction}
                <div class="flex justify-end gap-4">
                    {#if cancelText}
                        <!-- Two-button confirmation style -->
                        <button
                                type="button"
                                class="px-4 py-2 text-sm text-slate-600 hover:text-slate-900"
                                on:click={doCancel}
                                disabled={isLoading}
                        >
                            {cancelText}
                        </button>
                        <button
                                type="button"
                                class="px-4 py-2 text-sm rounded-md {isDangerous ? 'bg-red-500 hover:bg-red-600' : 'bg-sky-400 hover:bg-sky-500'} text-white"
                                on:click={doConfirm}
                                disabled={isLoading}
                        >
                            {confirmText}
                        </button>
                    {:else}
                        <!-- Single-button 'OK' style -->
                        <button
                                type="button"
                                class="px-4 py-2 text-sm bg-sky-400 hover:bg-sky-500 text-white rounded-md"
                                on:click={doConfirm}
                                disabled={isLoading}
                        >
                            {confirmText}
                        </button>
                    {/if}
                </div>
            {/if}
        </div>
    {/if}
</dialog>

<style lang="postcss">
    dialog[open] {
        opacity: 1;
        transform: translateY(0);
        transition: opacity 0.15s ease-out, transform 0.15s ease-out;
    }

    dialog {
        opacity: 0;
        transform: translateY(10px);
        margin: auto;
    }

    dialog::backdrop {
        opacity: 1;
        transition: opacity 0.15s ease-out, backdrop-filter 0.15s ease-out;
    }

    dialog:not([open])::backdrop {
        opacity: 0;
        backdrop-filter: blur(0px);
    }

    .loader {
        width: 20px;
        height: 20px;
        border: 3px solid #37393A;
        border-bottom-color: transparent;
        border-radius: 50%;
        display: inline-block;
        box-sizing: border-box;
        animation: rotation 1s linear infinite;
    }

    @keyframes rotation {
        0% {
            transform: rotate(0deg);
        }
        100% {
            transform: rotate(360deg);
        }
    }

</style>