<script lang="ts">
  import { invalidateAll } from "$app/navigation";
  import { getApiClient, handleApiError } from "$lib";
  import { Button } from "@nanoteck137/nano-ui";
  import toast from "svelte-5-french-toast";

  let { data } = $props();
  const apiClient = getApiClient();
</script>

<p>{data.user?.username}</p>

<Button href="/account/password">Change Password</Button>
<Button href="/account/tokens">Api Tokens</Button>

<Button
  onclick={async () => {
    const res = await apiClient.importMalAnimeList("nanoteck137");
    if (!res.success) {
      return handleApiError(res.error);
    }

    toast.success("Successfully imported mal list");
    invalidateAll();
  }}
>
  Test Import
</Button>
