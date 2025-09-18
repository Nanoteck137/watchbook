<script lang="ts">
  import { invalidateAll } from "$app/navigation";
  import { getApiClient, handleApiError } from "$lib";
  import { Button } from "@nanoteck137/nano-ui";
  import toast from "svelte-5-french-toast";

  const { data } = $props();
  const apiClient = getApiClient();
</script>

{#each data.parts as part}
  <div class="flex items-center gap-2">
    <p>{part.index}. {part.name}</p>
    <Button
      variant="link"
      onclick={async () => {
        const res = await apiClient.removePart(
          data.media.id,
          part.index.toString(),
        );
        if (!res.success) {
          return handleApiError(res.error);
        }

        toast.success("Successfully removed part");
        invalidateAll();
      }}
    >
      Remove
    </Button>
  </div>
{/each}
