<script lang="ts">
  import Spacer from "$lib/components/Spacer.svelte";
  import { Button } from "@nanoteck137/nano-ui";

  const { data, children } = $props();
</script>

<p class="text-center">Media Settings</p>

<Spacer size="md" />

<nav class="flex justify-center">
  <Button
    href="/media/{data.media.id}/settings"
    variant="link"
    data-sveltekit-noscroll>Update</Button
  >
  <Button
    href="/media/{data.media.id}/settings/images"
    variant="link"
    data-sveltekit-noscroll
  >
    Images
  </Button>
  <Button
    href="/media/{data.media.id}/settings/edit"
    variant="link"
    data-sveltekit-noscroll
  >
    Edit
  </Button>
  <Button
    href="/media/{data.media.id}/settings/parts"
    variant="link"
    data-sveltekit-noscroll
  >
    Parts
  </Button>
</nav>

{@render children()}
