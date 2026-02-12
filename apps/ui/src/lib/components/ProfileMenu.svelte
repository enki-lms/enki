<script lang="ts">
  import Avatar from "./Avatar.svelte";

  export let src: string | undefined = undefined;
  export let name: string = "";
  export let email: string = "";
  export let role: string = "";
  export let width: string = "80";
  export let height: string = "80";

  let isOpen = false;

  function toggleMenu() {
    isOpen = !isOpen;
  }

  function closeMenu() {
    isOpen = false;
  }

  function handleSignOut() {
    document.cookie = "token=; expires=Thu, 01 Jan 1970 00:00:00 UTC; path=/;";
    window.location.href = "/login";
  }

  function handleClickOutside(event: MouseEvent) {
    const target = event.target as HTMLElement;
    if (!target.closest(".profile-menu-container")) {
      closeMenu();
    }
  }

  function formatRole(r: string): string {
    if (!r) return "";
    return r.charAt(0).toUpperCase() + r.slice(1);
  }
</script>

<svelte:window on:click={handleClickOutside} />

<div class="profile-menu-container relative">
  <button
    type="button"
    class="cursor-pointer rounded-full transition-transform hover:scale-105 focus:outline-none focus:ring-2 focus:ring-blue-500 focus:ring-offset-2"
    on:click|stopPropagation={toggleMenu}
    aria-haspopup="true"
    aria-expanded={isOpen}
  >
    <Avatar {src} {name} {width} {height} />
  </button>

  {#if isOpen}
    <div
      class="absolute left-0 top-full mt-2 w-64 rounded-lg bg-white shadow-lg ring-1 ring-black ring-opacity-5 z-50"
      role="menu"
    >
      <div class="px-4 py-3 border-b border-gray-100">
        <p class="text-sm font-semibold text-gray-900 truncate">{name}</p>
        {#if email}
          <p class="text-xs text-gray-500 truncate">{email}</p>
        {/if}
        {#if role}
          <span
            class="inline-block mt-1 px-2 py-0.5 text-xs font-medium rounded-full bg-blue-100 text-blue-700"
          >
            {formatRole(role)}
          </span>
        {/if}
      </div>

      <a
        href="/submissions"
        class="flex w-full items-center gap-3 px-4 py-2 text-left text-sm text-gray-700 hover:bg-gray-100 transition-colors"
        role="menuitem"
        on:click={closeMenu}
      >
        <svg
          xmlns="http://www.w3.org/2000/svg"
          class="h-5 w-5 text-gray-500"
          fill="none"
          viewBox="0 0 24 24"
          stroke="currentColor"
          stroke-width="2"
        >
          <path
            stroke-linecap="round"
            stroke-linejoin="round"
            d="M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z"
          />
        </svg>
        Submission History
      </a>

      <button
        type="button"
        class="flex w-full items-center gap-3 px-4 py-2 text-left text-sm text-gray-700 hover:bg-gray-100 transition-colors"
        on:click={handleSignOut}
        role="menuitem"
      >
        <svg
          xmlns="http://www.w3.org/2000/svg"
          class="h-5 w-5 text-gray-500"
          fill="none"
          viewBox="0 0 24 24"
          stroke="currentColor"
          stroke-width="2"
        >
          <path
            stroke-linecap="round"
            stroke-linejoin="round"
            d="M17 16l4-4m0 0l-4-4m4 4H7m6 4v1a3 3 0 01-3 3H6a3 3 0 01-3-3V7a3 3 0 013-3h4a3 3 0 013 3v1"
          />
        </svg>
        Sign out
      </button>
    </div>
  {/if}
</div>
