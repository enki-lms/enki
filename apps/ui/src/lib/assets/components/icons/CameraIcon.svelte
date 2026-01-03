<script>
    // Add any interactivity here if needed
</script>

<style>
    .camera {
        width: 200px;
        height: 150px;
    }
</style>

<svg
    class="camera"
    viewBox="0 0 200 150"
    xmlns="http://www.w3.org/2000/svg"
>
    <!-- Camera Body -->
    <rect
        x="10"
        y="20"
        rx="15"
        ry="15"
        width="180"
        height="120"
        fill="#333"
        stroke="black"
        stroke-width="2"
    />

    <!-- Lens -->
    <circle
        cx="100"
        cy="80"
        r="35"
        fill="#222"
        stroke="black"
        stroke-width="2"
    />
    <circle
        cx="100"
        cy="80"
        r="25"
        fill="#444"
        stroke="black"
        stroke-width="2"
    />

    <!-- Flash -->
    <circle
        cx="160"
        cy="40"
        r="10"
        fill="#ccc"
        stroke="black"
        stroke-width="2"
    />

    <!-- Shutter Button -->
    <rect
        x="20"
        y="30"
        width="30"
        height="10"
        rx="5"
        ry="5"
        fill="#666"
        stroke="black"
        stroke-width="2"
    />
</svg>
