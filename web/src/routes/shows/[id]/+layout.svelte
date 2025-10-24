<script lang="ts">
  import BannerHeader from "$lib/components/BannerHeader.svelte";
  import Spacer from "$lib/components/Spacer.svelte";
  import { Breadcrumb, Button } from "@nanoteck137/nano-ui";
  import { isRoleAdmin } from "$lib/utils";

  const { data, children } = $props();
</script>

<Breadcrumb.Root class="py-2">
  <Breadcrumb.List>
    <Breadcrumb.Item>
      <Breadcrumb.Link href="/collections">Collections</Breadcrumb.Link>
    </Breadcrumb.Item>
    <Breadcrumb.Separator />
    <Breadcrumb.Item>
      <Breadcrumb.Page class="line-clamp-1 max-w-96 text-ellipsis">
        {data.show.name}
      </Breadcrumb.Page>
    </Breadcrumb.Item>
  </Breadcrumb.List>
</Breadcrumb.Root>

<Spacer size="md" />

<BannerHeader
  title={data.show.name}
  coverUrl={data.show.coverUrl}
  bannerUrl={data.show.bannerUrl}
  logoUrl={data.show.logoUrl}
>
  {#snippet buttons()}
    <div class="flex gap-4">
      <Button
        href="/shows/{data.show.id}"
        variant="secondary"
        data-sveltekit-noscroll>Overview</Button
      >

      {#if isRoleAdmin(data.user?.role)}
        <Button
          href="/shows/{data.show.id}/settings"
          variant="secondary"
          data-sveltekit-noscroll
        >
          Settings
        </Button>
      {/if}
    </div>
  {/snippet}
</BannerHeader>

<Spacer size="md" />

{@render children()}
