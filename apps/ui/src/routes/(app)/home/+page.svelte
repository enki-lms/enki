<script lang="ts">
    import Avatar from "$lib/components/Avatar.svelte";
    import LogoPlaceHolder from "$lib/components/LogoPlaceHolder.svelte";
    import TaskCard from "$lib/components/TaskCard.svelte";
</script>

<div class="min-h-screen bg-gray-200">
    <header class="p-2 flex">
        <div class="p-2"><Avatar width="80" height="80" src="src/lib/assets/images/temp-pfp.png"/></div>
 	    <div class="p-3"><LogoPlaceHolder /></div>
    </header>
	<main class="p-8">
		<div class="bg-white rounded-xl p-12 shadow-sm">
			<div class="grid grid-cols-2 gap-6">
				<!-- Ovdje ce ic zadaci -->
                
			</div>
		</div>
	</main>
</div>