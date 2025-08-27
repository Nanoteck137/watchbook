<script lang="ts">
  const { data } = $props();
</script>

{#if data.notificaitons}
  {#each data.notifications! as notification}
    <p>{notification.title}</p>
  {/each}
{/if}
