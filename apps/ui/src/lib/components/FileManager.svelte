<script lang="ts">
    import Button from '$lib/components/Button.svelte';
    import FolderIcon from '$lib/components/icons/FolderIcon.svelte';
    import FileIcon from '$lib/components/icons/FileIcon.svelte';
    import UploadIcon from '$lib/components/icons/UploadIcon.svelte';
    import FolderPlusIcon from '$lib/components/icons/FolderPlusIcon.svelte';
    import EmptyFolderIcon from '$lib/components/icons/EmptyFolderIcon.svelte';
    import ShareIcon from '$lib/components/icons/ShareIcon.svelte';
    import SortIcon from '$lib/components/icons/SortIcon.svelte';

    export let courseId: string | null = null;

    interface FileItem {
        name: string;
        modified: string;
        size: string;
        sharing: string;
        type: 'folder' | 'file';
        fileType?: 'document' | 'image' | 'zip';
        isShared?: boolean;
    }

    // za backend
    const files: FileItem[] = [];

    let sortBy: 'name' | 'modified' | 'size' | 'sharing' = 'name';
    let sortOrder: 'asc' | 'desc' = 'asc';

    const handleSort = (column: typeof sortBy) => {
        if (sortBy === column) {
            sortOrder = sortOrder === 'asc' ? 'desc' : 'asc';
        } else {
            sortBy = column;
            sortOrder = 'asc';
        }
    };

    const handleUpload = () => {
        // za backend
        console.log('Upload clicked for course:', courseId);
    };

    const handleNewFolder = () => {
        // za backend
        console.log('New folder clicked for course:', courseId);
    };
</script>

<div class="w-full">
    <div class="flex items-center justify-between mb-6">
        <h2 class="text-2xl font-semibold text-gray-900">My Files</h2>
        <div class="flex gap-3">
            <Button size="md" on:click={handleNewFolder}>
                <span slot="icon" class="text-white">
                    <FolderPlusIcon />
                </span>
                New Folder
            </Button>
            <Button size="md" on:click={handleUpload}>
                <span slot="icon" class="text-white">
                    <UploadIcon />
                </span>
                Upload
            </Button>
        </div>
    </div>

    <div class="bg-white rounded-lg border border-gray-200 overflow-hidden">
        <div class="flex items-center px-4 py-3 bg-gray-50 border-b border-gray-200 text-sm font-medium text-gray-700">
            <div class="flex items-center flex-1">
                <button 
                    class="flex items-center gap-1 hover:text-gray-900 transition-colors"
                    on:click={() => handleSort('name')}
                >
                    Name
                    {#if sortBy === 'name'}
                        <SortIcon direction={sortOrder} />
                    {/if}
                </button>
            </div>
            <div class="flex items-center gap-8">
                <button 
                    class="w-32 text-left hover:text-gray-900 transition-colors"
                    on:click={() => handleSort('modified')}
                >
                    Modified
                </button>
                <button 
                    class="w-24 text-left hover:text-gray-900 transition-colors"
                    on:click={() => handleSort('size')}
                >
                    File Size
                </button>
                <button 
                    class="w-32 text-left hover:text-gray-900 transition-colors"
                    on:click={() => handleSort('sharing')}
                >
                    Sharing
                </button>
            </div>
        </div>

        {#if files.length === 0}
            <div class="py-16 text-center text-gray-500">
                <EmptyFolderIcon />
                <p class="text-lg font-medium">No files yet</p>
                <p class="text-sm mt-1">Upload files or create a new folder to get started</p>
            </div>
        {:else}
            <div class="divide-y divide-gray-100">
                {#each files as file}
                    <button 
                        class="w-full flex items-center px-4 py-3 hover:bg-gray-50 transition-colors duration-150 border-b border-gray-100 group"
                    >
                        <div class="flex items-center flex-1 min-w-0">
                            <div class="flex-shrink-0 mr-3">
                                {#if file.type === 'folder'}
                                    <FolderIcon isShared={file.isShared || false} />
                                {:else}
                                    <FileIcon type={file.fileType || 'document'} />
                                {/if}
                            </div>
                            <div class="flex-1 min-w-0 text-left">
                                <p class="text-sm font-medium text-gray-900 truncate group-hover:text-gray-700">
                                    {file.name}
                                </p>
                            </div>
                        </div>
                        
                        <div class="flex items-center gap-8 text-sm text-gray-600">
                            <span class="w-32 text-left">{file.modified}</span>
                            <span class="w-24 text-left">{file.size}</span>
                            <span class="w-32 text-left flex items-center gap-1">
                                {#if file.isShared}
                                    <ShareIcon />
                                {/if}
                                {file.sharing}
                            </span>
                        </div>
                    </button>
                {/each}
            </div>
        {/if}
    </div>
</div>
