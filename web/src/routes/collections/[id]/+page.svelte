<script lang="ts">
  const { data } = $props();
</script>

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
