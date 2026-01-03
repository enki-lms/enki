<script lang="ts">
    export let color = "black";
    export let width = "26";
    export let height = "26";
</script>

<svg {width} {height} viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg">
    <path d="M9 18L15 12L9 6" fill={color}/>
    <path d="M9 18L15 12L9 6" stroke={color} stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/>
</svg>