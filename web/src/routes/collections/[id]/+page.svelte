<script lang="ts">
  const { data } = $props();
</script>

<!-- svelte-ignore a11y_missing_attribute -->
<div class="relative">
  <img
    class="h-[40vh] w-full object-cover object-top"
    src={data.collection.bannerUrl}
  />

  <!-- <div>

  </div> -->
  <!-- <img
    class="absolute right-[25vw] top-[10vh] z-10 h-[16vh]"
    src={data.collection.logoUrl}
  />
  <div class="absolute inset-0 bg-gray-600/75"></div>

  <div class="absolute bottom-0 h-16 w-full bg-gray-900"></div>
  <img
    class="absolute bottom-4 left-4 h-80 w-56"
    src={data.collection.coverUrl}
  />

  <div class="absolute bottom-0 left-44 right-0 h-16 bg-red-200"></div> -->
</div>

<p>{data.collection.name}</p>

<div class="flex flex-col">
  {#each data.groups as group}
    <p>{group.name}</p>
    <div class="m-2 flex flex-col bg-red-200">
      {#each group.entries as entry}
        <p>{entry.collectionName}</p>
      {/each}
    </div>
  {/each}
</div>
