<script lang="ts">
  import { Button } from "@nanoteck137/nano-ui";
  import ProviderSearchModal from "./ProviderSearchModal.svelte";
  import { getApiClient, handleApiError } from "$lib";
  import toast from "svelte-5-french-toast";
  import { invalidateAll } from "$app/navigation";
  import Provider from "./Provider.svelte";

  const { data } = $props();
  const apiClient = getApiClient();
</script>

<div class="flex flex-col gap-2">
  {#each data.providers as provider}
    <Provider {provider} />
  {/each}
</div>
