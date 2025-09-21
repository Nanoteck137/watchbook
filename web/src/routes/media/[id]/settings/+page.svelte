<script>
  import { getApiClient, handleApiError } from "$lib";
  import ConfirmBox from "$lib/components/ConfirmBox.svelte";
  import toast from "svelte-5-french-toast";
  import ProviderUpdateButton from "../ProviderUpdateButton.svelte";
  import { goto } from "$app/navigation";
  import { Button } from "@nanoteck137/nano-ui";
  import { Trash } from "lucide-svelte";

  const { data } = $props();
  const apiClient = getApiClient();

  let openDeleteModal = $state(false);
</script>

<div class="mx-auto max-w-3xl p-6">
  <!-- Media Settings -->
  <div class="space-y-8 rounded-2xl bg-white p-6 shadow">
    <!-- Header -->
    <header>
      <h1 class="text-2xl font-bold text-gray-900">Media Settings</h1>
      <p class="text-sm text-gray-500">
        Admin tools for editing and maintaining this media entry.
      </p>
    </header>

    <!-- Edit Media Info -->
    <section>
      <h2 class="mb-3 text-lg font-semibold text-gray-900">Edit Media</h2>
      <form class="space-y-4">
        <div>
          <label class="block text-sm font-medium text-gray-700">Title</label>
          <input
            type="text"
            value="Attack on Titan"
            class="mt-1 block w-full rounded-lg border-gray-300 shadow-sm focus:border-indigo-500 focus:ring-indigo-500"
          />
        </div>

        <div>
          <label class="block text-sm font-medium text-gray-700"
            >Description</label
          >
          <textarea
            rows="4"
            class="mt-1 block w-full rounded-lg border-gray-300 shadow-sm focus:border-indigo-500 focus:ring-indigo-500"
            >A story about humanity fighting titans...</textarea
          >
        </div>

        <div class="flex gap-4">
          <div class="flex-1">
            <label class="block text-sm font-medium text-gray-700"
              >Start Date</label
            >
            <input
              type="date"
              class="mt-1 block w-full rounded-lg border-gray-300 shadow-sm focus:border-indigo-500 focus:ring-indigo-500"
            />
          </div>
          <div class="flex-1">
            <label class="block text-sm font-medium text-gray-700"
              >End Date</label
            >
            <input
              type="date"
              class="mt-1 block w-full rounded-lg border-gray-300 shadow-sm focus:border-indigo-500 focus:ring-indigo-500"
            />
          </div>
        </div>

        <div>
          <label class="block text-sm font-medium text-gray-700"
            >Tags (comma-separated)</label
          >
          <input
            type="text"
            value="Action, Drama, Fantasy"
            class="mt-1 block w-full rounded-lg border-gray-300 shadow-sm focus:border-indigo-500 focus:ring-indigo-500"
          />
        </div>

        <button
          type="submit"
          class="rounded-lg bg-indigo-600 px-4 py-2 text-white hover:bg-indigo-700"
        >
          Save Changes
        </button>
      </form>
    </section>

    <!-- Maintenance Actions -->
    <section>
      <h2 class="mb-3 text-lg font-semibold text-gray-900">Maintenance</h2>
      <div class="flex flex-wrap gap-3">
        <button
          class="rounded-lg bg-blue-600 px-4 py-2 text-white hover:bg-blue-700"
        >
          Send Update Request
        </button>
        <button
          class="rounded-lg bg-yellow-500 px-4 py-2 text-white hover:bg-yellow-600"
        >
          Refresh from API
        </button>
        <button
          class="rounded-lg bg-red-600 px-4 py-2 text-white hover:bg-red-700"
        >
          Delete Media
        </button>
      </div>
    </section>
  </div>
</div>

<Button
  variant="destructive"
  onclick={() => {
    openDeleteModal = true;
  }}
>
  <Trash />
  Delete Media
</Button>

{#each data.media.providers as provider}
  <ProviderUpdateButton mediaId={data.media.id} {provider} />
{/each}

<ConfirmBox
  bind:open={openDeleteModal}
  title="Delete Media?"
  description="Are you sure you want to delete this media? This action cannot be undone."
  confirmText="Delete"
  onResult={async () => {
    const res = await apiClient.deleteMedia(data.media.id);
    if (!res.success) {
      return handleApiError(res.error);
    }

    toast.success("Successfully deleted media");
    goto(`/media`, { invalidateAll: true });
  }}
/>
