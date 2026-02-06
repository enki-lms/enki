<script lang="ts">
  import { goto } from "$app/navigation";
  import { onMount } from "svelte";

  const hasToken = () =>
    document.cookie.split(";").some((c) => c.trim().startsWith("token="));

  onMount(() => {
    goto(hasToken() ? "/home" : "/login");
  });
</script>
