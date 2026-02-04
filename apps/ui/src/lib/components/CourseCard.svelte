<script lang="ts">
    export let name: string;
    export let description: string;
    export let color: string = 'blue';

    const colorClasses: Record<string, string> = {
        blue: 'bg-blue-50 border-blue-200 hover:bg-blue-100',
        purple: 'bg-purple-50 border-purple-200 hover:bg-purple-100',
        green: 'bg-green-50 border-green-200 hover:bg-green-100',
        orange: 'bg-orange-50 border-orange-200 hover:bg-orange-100',
        red: 'bg-red-50 border-red-200 hover:bg-red-100',
        yellow: 'bg-yellow-50 border-yellow-200 hover:bg-yellow-100',
    };

    const iconColorClasses: Record<string, string> = {
        blue: 'text-blue-600',
        purple: 'text-purple-600',
        green: 'text-green-600',
        orange: 'text-orange-600',
        red: 'text-red-600',
        yellow: 'text-yellow-600',
    };
</script>

<button
    class="w-full p-6 rounded-xl border-2 transition-all duration-200 text-left {colorClasses[color] || colorClasses.blue} hover:shadow-md"
    on:click
>
    <div class="flex items-start gap-4">
        <div class="flex-shrink-0 mt-1">
            <svg class="w-8 h-8 {iconColorClasses[color] || iconColorClasses.blue}" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                <path d="M2 3h6a4 4 0 0 1 4 4v14a3 3 0 0 0-3-3H2z"></path>
                <path d="M22 3h-6a4 4 0 0 0-4 4v14a3 3 0 0 1 3-3h7z"></path>
            </svg>
        </div>
        <div class="flex-1 min-w-0">
            <h3 class="text-lg font-semibold text-gray-900 mb-1 truncate">{name}</h3>
            <p class="text-sm text-gray-600 line-clamp-2">{description}</p>
        </div>
        <div class="flex-shrink-0">
            <svg class="w-5 h-5 text-gray-400" viewBox="0 0 20 20" fill="currentColor">
                <path fill-rule="evenodd" d="M7.293 14.707a1 1 0 010-1.414L10.586 10 7.293 6.707a1 1 0 011.414-1.414l4 4a1 1 0 010 1.414l-4 4a1 1 0 01-1.414 0z" clip-rule="evenodd" />
            </svg>
        </div>
    </div>
</button>
