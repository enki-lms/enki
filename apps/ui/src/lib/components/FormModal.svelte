<script lang="ts" module>
  export interface FieldConfig {
    name: string;
    label: string;

    type: "text" | "textarea" | "number" | "select" | "file";
    placeholder?: string;
    required?: boolean;
    options?: { value: string; label: string }[];
    min?: number;
    max?: number;
    accept?: string;
  }
</script>

<script lang="ts">
  import { createEventDispatcher } from "svelte";
  import { fade, scale } from "svelte/transition";
  import Button from "$lib/components/Button.svelte";

  export let isOpen = false;
  export let title: string;
  export let submitText = "Create";
  export let isLoading = false;

  export let fields: FieldConfig[] = [];
  export let initialValues: Record<string, any> = {};

  let values: Record<string, any> = {};
  let prevIsOpen = false;

  $: if (isOpen && !prevIsOpen) {
    // Only reset values when modal first opens
    values = { ...initialValues };
    fields.forEach((field) => {
      if (values[field.name] === undefined) {
        values[field.name] = field.type === "number" ? 0 : "";
      }
    });
    prevIsOpen = true;
  } else if (!isOpen && prevIsOpen) {
    prevIsOpen = false;
  }

  const dispatch = createEventDispatcher<{
    submit: Record<string, any>;
    cancel: void;
  }>();

  const handleSubmit = () => {
    const isValid = fields.every((field) => {
      if (!field.required) return true;
      const value = values[field.name];
      if (typeof value === "string") return value.trim() !== "";
      return value !== undefined && value !== null;
    });

    if (isValid) {
      dispatch("submit", { ...values });
    }
  };

  const handleCancel = () => {
    dispatch("cancel");
    isOpen = false;
  };

  const handleKeydown = (e: KeyboardEvent) => {
    if (e.key === "Escape") {
      handleCancel();
    }
  };
</script>

{#if isOpen}
  <div
    class="fixed inset-0 bg-black/30 flex items-center justify-center z-50 p-4"
    on:click={handleCancel}
    on:keydown={handleKeydown}
    role="button"
    tabindex="0"
    transition:fade={{ duration: 200 }}
  >
    <div
      class="bg-white rounded-2xl shadow-xl max-w-lg w-full p-6 max-h-[90vh] overflow-y-auto"
      on:click|stopPropagation
      role="dialog"
      aria-modal="true"
      transition:scale={{ duration: 200, start: 0.95 }}
    >
      <h2 class="text-2xl font-semibold text-gray-900 mb-6">{title}</h2>

      <form on:submit|preventDefault={handleSubmit} class="space-y-4">
        {#each fields as field}
          <div class="space-y-1">
            <label
              for={field.name}
              class="block text-sm font-medium text-gray-700"
            >
              {field.label}
              {#if field.required}<span class="text-red-500">*</span>{/if}
            </label>

            {#if field.type === "text"}
              <input
                id={field.name}
                type="text"
                value={values[field.name] || ""}
                on:input={(e) => {
                  values[field.name] = e.currentTarget.value;
                  values = { ...values };
                }}
                placeholder={field.placeholder}
                class="w-full px-4 py-3 border-2 border-gray-200 rounded-lg focus:outline-none focus:border-sky-400 transition-colors"
                required={field.required}
                disabled={isLoading}
              />
            {:else if field.type === "textarea"}
              <textarea
                id={field.name}
                value={values[field.name] || ""}
                on:input={(e) => {
                  values[field.name] = e.currentTarget.value;
                  values = { ...values };
                }}
                placeholder={field.placeholder}
                rows="4"
                class="w-full px-4 py-3 border-2 border-gray-200 rounded-lg focus:outline-none focus:border-sky-400 transition-colors resize-none"
                required={field.required}
                disabled={isLoading}
              ></textarea>
            {:else if field.type === "number"}
              <input
                id={field.name}
                type="number"
                value={values[field.name] ?? ""}
                on:input={(e) => {
                  const val = e.currentTarget.value;
                  values[field.name] = val === "" ? 0 : Number(val);
                  values = { ...values };
                }}
                placeholder={field.placeholder}
                min={field.min}
                max={field.max}
                class="w-full px-4 py-3 border-2 border-gray-200 rounded-lg focus:outline-none focus:border-sky-400 transition-colors"
                required={field.required}
                disabled={isLoading}
              />
            {:else if field.type === "select"}
              <select
                id={field.name}
                value={values[field.name] || ""}
                on:change={(e) => {
                  values[field.name] = e.currentTarget.value;
                  values = { ...values };
                }}
                class="w-full px-4 py-3 border-2 border-gray-200 rounded-lg focus:outline-none focus:border-sky-400 transition-colors bg-white"
                required={field.required}
                disabled={isLoading}
              >
                <option value="">{field.placeholder || "Select..."}</option>
                {#if field.options}
                  {#each field.options as option}
                    <option value={option.value}>{option.label}</option>
                  {/each}
                {/if}
              </select>
            {:else if field.type === "file"}
              <div class="flex items-center gap-2">
                <input
                  id={field.name}
                  type="file"
                  accept={field.accept}
                  on:change={(e) => {
                    const files = e.currentTarget.files;
                    if (files && files.length > 0) {
                      values[field.name] = files[0];
                      values = { ...values };
                    }
                  }}
                  class="w-full px-4 py-3 border-2 border-gray-200 rounded-lg focus:outline-none focus:border-sky-400 transition-colors bg-white file:mr-4 file:py-2 file:px-4 file:rounded-full file:border-0 file:text-sm file:font-semibold file:bg-sky-50 file:text-sky-700 hover:file:bg-sky-100"
                  required={field.required && !values[field.name]}
                  disabled={isLoading}
                />
                {#if values[field.name] && typeof values[field.name] === "string"}
                  <span class="text-xs text-green-600 truncate max-w-[150px]"
                    >{values[field.name].split("/").pop()}</span
                  >
                {/if}
              </div>
              {#if values[field.name] && typeof values[field.name] === "string" && values[field.name].startsWith("http")}
                <img
                  src={values[field.name]}
                  alt="Preview"
                  class="mt-2 h-20 w-auto object-contain border rounded"
                />
              {/if}
            {/if}
          </div>
        {/each}

        <div class="flex gap-3 mt-6 justify-end pt-4">
          <button
            type="button"
            on:click={handleCancel}
            class="px-6 py-2 rounded-full border-2 border-gray-300 text-gray-700 font-semibold hover:bg-gray-50 transition-colors"
            disabled={isLoading}
          >
            Cancel
          </button>
          <button
            type="button"
            on:click={handleSubmit}
            disabled={isLoading}
            class="px-6 py-2 rounded-full bg-sky-400 hover:bg-sky-500 text-white font-semibold transition-colors disabled:opacity-50 disabled:cursor-not-allowed"
          >
            {#if isLoading}
              <span
                class="inline-block w-4 h-4 border-2 border-white border-t-transparent rounded-full animate-spin mr-2"
              ></span>
            {/if}
            {submitText}
          </button>
        </div>
      </form>
    </div>
  </div>
{/if}
