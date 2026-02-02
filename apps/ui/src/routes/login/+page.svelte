<script lang="ts">
  import Button from "$lib/components/Button.svelte";
  import Divider from "$lib/components/Divider.svelte";
  import enkiLogo from "$lib/assets/images/enki.png"
  import "./bg.css";
</script>

<div class="magicpattern min-h-screen flex items-center justify-center px-4">
  <div class="bg-white rounded-2xl shadow-xl px-10 py-8 w-full max-w-md">
    <div class="flex justify-center mb-4">
      <img src={enkiLogo} alt="Enki" class="h-20 w-20 rounded-full object-cover" />
    </div>

    <h1 class="text-center text-2xl font-semibold text-gray-900 mb-2">
      Welcome to Enki
    </h1>

    <div class="mx-auto mb-6 h-0.5 w-12 bg-blue-500 rounded"></div>

    <a
      href="http://localhost:8080/auth/login"
      data-sveltekit-reload
      class="flex w-full justify-center rounded-xl py-3 text-sm font-medium bg-blue-400 hover:bg-blue-500 text-white shadow-md transition-colors"
    >
      Sign in with AAI@EduHr
    </a>

    <div class="my-6">
      <Divider />
    </div>

    <p class="text-center text-xs text-gray-500 leading-relaxed">
      By continuing, you agree to our
      <a href="#" class="text-blue-500 hover:underline">Terms of Service</a>
      and
      <a href="#" class="text-blue-500 hover:underline">Privacy Policy</a>.
    </p>
  </div>
</div>

<style>
  :global(body) {
    margin: 0;
    padding: 0;
    background-color: #f9fafb;
  }
</style>
