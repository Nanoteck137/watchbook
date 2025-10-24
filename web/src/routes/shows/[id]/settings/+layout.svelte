<script lang="ts">
  import { Card, ScrollArea } from "@nanoteck137/nano-ui";

  const { data, children } = $props();
</script>

<Card.Root>
  <div class="flex w-full justify-center gap-4 border-b">
    <ScrollArea class="px-2 pb-4 sm:pb-0" orientation="horizontal">
      <div class="flex gap-4 pt-4">
        <a
          class="border-b-2 px-2 py-3 text-sm font-medium hover:brightness-75"
          href="/collections/{data.collection.id}/settings"
          data-sveltekit-noscroll
        >
          General
        </a>

        <a
          class="border-b-2 px-2 py-3 text-sm font-medium hover:brightness-75"
          href="/collections/{data.collection.id}/settings/images"
          data-sveltekit-noscroll
        >
          Images
        </a>

        <a
          class="border-b-2 px-2 py-3 text-sm font-medium hover:brightness-75"
          href="/collections/{data.collection.id}/settings/edit"
          data-sveltekit-noscroll
        >
          Edit
        </a>
      </div>
    </ScrollArea>
  </div>

  <div class="p-6">
    {@render children()}
  </div>
</Card.Root>
