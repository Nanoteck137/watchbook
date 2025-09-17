<script lang="ts">
  import CollectionCard from "$lib/components/CollectionCard.svelte";
  import Spacer from "$lib/components/Spacer.svelte";
  import { Button } from "@nanoteck137/nano-ui";
  import NewCollectionModal from "./NewCollectionModal.svelte";
  import Filter from "./Filter.svelte";
  import StandardPagination from "$lib/components/StandardPagination.svelte";

  const { data } = $props();

  let openNewCollectionModal = $state(false);
</script>

<Filter fullFilter={data.filter} />

<p>Total collections: {data.page.totalItems}</p>

<Button
  onclick={() => {
    openNewCollectionModal = true;
  }}
>
  New Collection
</Button>

<div class="mx-auto max-w-7xl px-4 sm:px-6 lg:px-8">
  <div
    class="grid grid-cols-[repeat(auto-fit,minmax(280px,1fr))] items-center justify-items-center gap-6"
  >
    {#each data.collections as collection}
      <CollectionCard
        href="/collections/{collection.id}"
        name={collection.name}
        coverUrl={collection.coverUrl}
      />
    {/each}
  </div>
</div>

<Spacer size="sm" />

<StandardPagination pageData={data.page} />

<NewCollectionModal bind:open={openNewCollectionModal} />
