<script lang="ts">
  const { data } = $props();
</script>

<p>{data.collection.name}</p>

<div class="flex flex-col">
  {#each data.items as order}
    <div class="m-2 flex flex-col bg-red-200">
      {#each order as item}
        <p>{item.collectionName}</p>
      {/each}
    </div>
  {/each}
</div>
