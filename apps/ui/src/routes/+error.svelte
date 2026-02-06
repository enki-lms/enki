<script lang="ts">
  import { page } from "$app/stores";
  import enkiLogo from "$lib/assets/images/enki.png";
  import "./login/bg.css";
</script>

<div class="magicpattern min-h-screen flex items-center justify-center px-4">
  <div
    class="bg-white rounded-2xl shadow-xl px-10 py-8 w-full max-w-md text-center"
  >
    <div class="flex justify-center mb-4">
      <img
        src={enkiLogo}
        alt="Enki"
        class="h-20 w-20 rounded-full object-cover"
      />
    </div>

    <h1 class="text-6xl font-bold text-blue-500 mb-2">
      {$page.status}
    </h1>

    <h2 class="text-xl font-semibold text-gray-900 mb-4">
      {#if $page.status === 404}
        Page Not Found
      {:else}
        Something Went Wrong
      {/if}
    </h2>

    <p class="text-gray-500 text-sm mb-6">
      {#if $page.error?.message}
        {$page.error.message}
      {:else if $page.status === 404}
        The page you're looking for doesn't exist or has been moved.
      {:else}
        An unexpected error occurred. Please try again later.
      {/if}
    </p>

    <a
      href="/home"
      class="inline-flex items-center justify-center rounded-xl py-3 px-6 text-sm font-medium bg-blue-400 hover:bg-blue-500 text-white shadow-md transition-colors"
    >
      Go to Home
    </a>
  </div>
</div>

<style>
  :global(body) {
    margin: 0;
    padding: 0;
    background-color: #f9fafb;
  }
</style>
