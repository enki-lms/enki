<script lang="ts">
    import { createEventDispatcher } from 'svelte';
    import EyeIcon from "$lib/assets/components/icons/EyeIcon.svelte";
    import EyeOffIcon from "$lib/assets/components/icons/EyeOffIcon.svelte";

    const dispatch = createEventDispatcher();

    export let placeholderText: string = "Enter password";
    export let value: string = "";

    let passwordVisible = false;
    let passwordInput: HTMLInputElement;

    function togglePasswordVisibility() {
        passwordVisible = !passwordVisible;
        passwordInput.type = passwordVisible ? 'text' : 'password';
    }

    function handleInput(event: Event) {
        value = (event.target as HTMLInputElement).value;
        dispatch('input', { value });
    }
</script>

<div class="relative w-full h-12">
    <input
            bind:this={passwordInput}
            type="password"
            placeholder={placeholderText}
            value={value}
            class="w-full h-full px-3 py-2 pr-12 border-b-2 text-lg leading-tight text-slate-900 bg-transparent border-slate-300 focus:outline-hidden focus:border-sky-400"
            on:input={handleInput}
    />
    <button
            type="button"
            class="absolute inset-y-0 right-0 flex items-center px-3 text-slate-400 hover:text-slate-600 transition-colors"
            on:click={togglePasswordVisibility}
            aria-label={passwordVisible ? "Hide password" : "Show password"}
    >
        {#if passwordVisible}
            <EyeOffIcon />
        {:else}
            <EyeIcon />
        {/if}
    </button>
</div>
