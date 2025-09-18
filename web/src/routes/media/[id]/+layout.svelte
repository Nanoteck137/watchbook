<script lang="ts">
  import BannerHeader from "$lib/components/BannerHeader.svelte";
  import Spacer from "$lib/components/Spacer.svelte";
  import { Breadcrumb, Button } from "@nanoteck137/nano-ui";
  import MediaDropdown from "./MediaDropdown.svelte";

  const { data, children } = $props();
</script>

<Breadcrumb.Root class="py-2">
  <Breadcrumb.List>
    <Breadcrumb.Item>
      <Breadcrumb.Link href="/media">Media</Breadcrumb.Link>
    </Breadcrumb.Item>
    <Breadcrumb.Separator />
    <Breadcrumb.Item>
      <Breadcrumb.Page class="line-clamp-1 max-w-96 text-ellipsis">
        {data.media.title}
      </Breadcrumb.Page>
    </Breadcrumb.Item>
  </Breadcrumb.List>
</Breadcrumb.Root>

<Spacer size="md" />

<BannerHeader
  title={data.media.title}
  description={data.media.description}
  coverUrl={data.media.coverUrl}
  bannerUrl={data.media.bannerUrl}
  logoUrl={data.media.logoUrl}
>
  {#snippet imageContent()}
    <MediaDropdown media={data.media} />
  {/snippet}

  {#snippet buttons()}
    <div class="flex gap-4">
      <Button
        href="/media/{data.media.id}"
        variant="secondary"
        data-sveltekit-noscroll
      >
        Overview
      </Button>
      <Button
        href="/media/{data.media.id}/parts"
        variant="secondary"
        data-sveltekit-noscroll
      >
        Parts
      </Button>
      <Button
        href="/media/{data.media.id}/settings"
        variant="secondary"
        data-sveltekit-noscroll
      >
        Settings
      </Button>

      <!-- <a class="bg-red-200" href="?">Overview</a>
        <a href="?">Parts</a> -->
      <!-- <Button variant="link">Overview</Button> -->
    </div>
  {/snippet}
</BannerHeader>

<Spacer size="md" />

{@render children()}
