<!--
<script lang="ts">
    export let type: string;
    export let placeholder: string;
    export let id: string;
</script>

<div class="relative">
    <input
            type={type}
            id={id}
            placeholder={placeholder}
            class="w-full py-2 pl-3 pr-12 text-sm leading-none text-gray-700 bg-transparent border-0 border-b-2 border-gray-300 focus:outline-hidden focus:ring-0 focus:border-blue-500"
    />
    <div class="absolute inset-y-0 right-0 flex items-center px-2 pointer-events-none">
        <svg class="w-4 h-4 text-gray-600" fill="currentColor" viewBox="0 0 20 20" xmlns="http://www.w3.org/2000/svg"><path fill-rule="evenodd" d="M10 18a8 8 0 100-16 8 8 0 000 16zm1-11a1 1 0 11-2 0 1 1 0 012 0zm-1 2a1 1 0 00-1 1v2a1 1 0 102 0v-2a1 1 0 00-1-1z" clip-rule="evenodd"></path></svg>
    </div>
</div>-->

<script lang="ts">
    import { createEventDispatcher } from 'svelte';
    import { onMount } from 'svelte';
    import { tick } from 'svelte';

    const dispatch = createEventDispatcher();

    export let isPassword: boolean = false;
    export let isCheckbox: boolean = false;
    export let isSearch: boolean = false;
    export let placeholderText: string = "Enter text";

    let textInput: HTMLInputElement;
    let checkboxVisible = false;
    let passwordVisible = false;

    onMount(() => {
        checkboxVisible = isCheckbox;
        passwordVisible = isPassword;
    });

    function togglePasswordVisibility() {
        passwordVisible = !passwordVisible;
        textInput.type = passwordVisible ? 'text' : 'password';
    }

    function handleInput(event: Event) {
        dispatch('textchange', { value: (event.target as HTMLInputElement).value });
    }
</script>

<div class="relative w-full h-12">
    <input
            bind:this={textInput}
            type={isPassword ? 'password' : 'text'}
            placeholder={placeholderText}
            class="w-full h-full px-3 py-2 border-b-2 text-lg leading-tight text-slate-900 bg-transparent border-slate-300 focus:outline-hidden focus:border-sky-400"
            on:input={handleInput}
    />
    {#if isSearch}
    <span class="absolute inset-y-0 left-0 flex items-center pl-2">
      <svg class="w-5 h-5 text-slate-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z" />
      </svg>
    </span>
    {/if}
    {#if isPassword}
        <button class="absolute inset-y-0 right-0 flex items-center px-3 text-slate-400" on:click={togglePasswordVisibility}>
            <svg class="w-5 h-5" fill="currentColor" viewBox="0 0 20 20">
                {#if passwordVisible}
                    <path fill-rule="evenodd" d="M10 2a8 8 0 018 8c0 1.493-.423 3.033-1.245 4.376A10.457 10.457 0 0110 18a10.457 10.457 0 01-6.755-3.624C1.423 13.033 1 11.493 1 10a8 8 0 019-8zm0 4a1 1 0 100 2 1 1 0 000-2z" clip-rule="evenodd"></path>
                {:else}
                    <path fill-rule="evenodd" d="M10 2a8 8 0 018 8c0 1.493-.423 3.033-1.245 4.376A10.457 10.457 0 0110 18a10.457 10.457 0 01-6.755-3.624C1.423 13.033 1 11.493 1 10a8 8 0 019-8zm0 12a4 4 0 100-8 4 4 0 000 8z" clip-rule="evenodd"></path>
                {/if}
            </svg>
        </button>
    {/if}
</div>
