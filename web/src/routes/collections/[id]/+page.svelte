<script lang="ts">
  const { data } = $props();
</script>

<p>{data.collection.name}</p>

<div class="flex flex-col">
  {#each data.items as item}
    <p>{item.collectionName}</p>
  {/each}
</div>
