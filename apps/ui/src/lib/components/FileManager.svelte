<script lang="ts">
    import Button from '$lib/components/Button.svelte';
    import FolderIcon from '$lib/components/icons/FolderIcon.svelte';
    import FileIcon from '$lib/components/icons/FileIcon.svelte';
    import UploadIcon from '$lib/components/icons/UploadIcon.svelte';
    import FolderPlusIcon from '$lib/components/icons/FolderPlusIcon.svelte';

    interface FileItem {
        name: string;
        modified: string;
        size: string;
        sharing: string;
        type: 'folder' | 'file';
        fileType?: 'document' | 'image' | 'zip';
        isShared?: boolean;
    }

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
        console.log('Upload clicked');
    };

    const handleNewFolder = () => {
        console.log('New folder clicked');
    };
</script>

<div class="w-full">
    <div class="flex items-center justify-between mb-6">
        <h2 class="text-2xl font-semibold text-gray-900">Course Files</h2>
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
                        <svg class="w-4 h-4 transform {sortOrder === 'desc' ? 'rotate-180' : ''}" viewBox="0 0 16 16" fill="currentColor">
                            <path d="M8 3l4 5H4l4-5z"/>
                        </svg>
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
                <svg class="w-16 h-16 mx-auto mb-4 text-gray-300" viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg">
                    <path d="M3 8C3 6.89543 3.89543 6 5 6H8.17157C8.70201 6 9.21071 5.78929 9.58579 5.41421L10.4142 4.58579C10.7893 4.21071 11.298 4 11.8284 4H19C20.1046 4 21 4.89543 21 6V18C21 19.1046 20.1046 20 19 20H5C3.89543 20 3 19.1046 3 18V8Z" stroke="currentColor" stroke-width="2"/>
                </svg>
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
                                    <svg class="w-4 h-4" viewBox="0 0 16 16" fill="none" xmlns="http://www.w3.org/2000/svg">
                                        <path d="M11 7C12.1046 7 13 6.10457 13 5C13 3.89543 12.1046 3 11 3C9.89543 3 9 3.89543 9 5C9 6.10457 9.89543 7 11 7Z" fill="#60A5FA"/>
                                        <path d="M5 10C6.10457 10 7 9.10457 7 8C7 6.89543 6.10457 6 5 6C3.89543 6 3 6.89543 3 8C3 9.10457 3.89543 10 5 10Z" fill="#60A5FA"/>
                                        <path d="M11 13C12.1046 13 13 12.1046 13 11C13 9.89543 12.1046 9 11 9C9.89543 9 9 9.89543 9 11C9 12.1046 9.89543 13 11 13Z" fill="#60A5FA"/>
                                        <path d="M7 8L9 6M7 8L9 10" stroke="#60A5FA" stroke-width="1.5"/>
                                    </svg>
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