<script lang="ts">
  import Errors from "$lib/components/Errors.svelte";
  import FormItem from "$lib/components/FormItem.svelte";
  import { Button, Input, Label, Select } from "@nanoteck137/nano-ui";
  import { zod } from "sveltekit-superforms/adapters";
  import { defaults, superForm } from "sveltekit-superforms/client";
  import {
    defaultSort,
    FullFilter,
    mediaRatings,
    mediaStatus,
    mediaTypes,
    sortTypes,
  } from "./types";
  import { goto } from "$app/navigation";
  import { page } from "$app/stores";

  export type Props = {
    fullFilter: FullFilter;
  };

  const { fullFilter }: Props = $props();

  function submit(data: FullFilter) {
    console.log(data);

    setTimeout(() => {
      const query = $page.url.searchParams;
      query.delete("query");
      query.delete("sort");

      query.delete("filterType");
      query.delete("filterStatus");
      query.delete("filterRating");

      query.delete("excludeType");
      query.delete("excludeStatus");
      query.delete("excludeRating");

      query.set("query", data.query);
      query.set("sort", data.sort);

      if (data.filter.type.length > 0) {
        query.set("filterType", data.filter.type.join(","));
      }

      if (data.filter.status.length > 0) {
        query.set("filterStatus", data.filter.status.join(","));
      }

      if (data.filter.rating.length > 0) {
        query.set("filterRating", data.filter.rating.join(","));
      }

      if (data.excludes.type.length > 0) {
        query.set("excludeType", data.excludes.type.join(","));
      }

      if (data.excludes.status.length > 0) {
        query.set("excludeStatus", data.excludes.status.join(","));
      }

      if (data.excludes.rating.length > 0) {
        query.set("excludeRating", data.excludes.rating.join(","));
      }

      // if (data.excludes.length > 0) {
      //   query.delete("filter");
      //   query.set("excludes", data.excludes.join(","));
      // }

      // if (data.filter !== "") {
      //   query.delete("excludes");
      //   query.set("filter", data.filter);
      // }

      goto("?" + query.toString(), { invalidateAll: true });
    }, 0);
  }

  const { form, errors, enhance, reset } = superForm(
    defaults(fullFilter, zod(FullFilter)),
    {
      id: "filter",
      SPA: true,
      validators: zod(FullFilter),
      dataType: "json",
      resetForm: false,
      onUpdate({ form }) {
        if (form.valid) {
          submit(form.data);
        }
      },
    },
  );
</script>

<form action="GET" class="flex flex-col gap-4" use:enhance>
  <FormItem>
    <Label for="query">Search</Label>
    <Input id="query" name="query" type="text" bind:value={$form.query} />
    <Errors errors={$errors.query} />
  </FormItem>

  <div class="flex flex-col gap-4">
    <div class="flex items-center gap-4">
      <p class="w-20">Filter</p>

      <Select.Root type="multiple" bind:value={$form.filter.type}>
        <Select.Trigger class="max-w-[120px]">Type</Select.Trigger>
        <Select.Content>
          {#each mediaTypes as ty (ty.value)}
            <Select.Item value={ty.value} label={ty.label} />
          {/each}
        </Select.Content>
      </Select.Root>

      <Select.Root type="multiple" bind:value={$form.filter.status}>
        <Select.Trigger class="max-w-[120px]">Status</Select.Trigger>
        <Select.Content>
          {#each mediaStatus as status (status.value)}
            <Select.Item value={status.value} label={status.label} />
          {/each}
        </Select.Content>
      </Select.Root>

      <Select.Root type="multiple" bind:value={$form.filter.rating}>
        <Select.Trigger class="max-w-[120px]">Rating</Select.Trigger>
        <Select.Content>
          {#each mediaRatings as rating (rating.value)}
            <Select.Item value={rating.value} label={rating.label} />
          {/each}
        </Select.Content>
      </Select.Root>
    </div>

    <div class="flex items-center gap-4">
      <p class="w-20">Excludes</p>

      <Select.Root type="multiple" bind:value={$form.excludes.type}>
        <Select.Trigger class="max-w-[120px]">Type</Select.Trigger>
        <Select.Content>
          {#each mediaTypes as ty (ty.value)}
            <Select.Item value={ty.value} label={ty.label} />
          {/each}
        </Select.Content>
      </Select.Root>

      <Select.Root type="multiple" bind:value={$form.excludes.status}>
        <Select.Trigger class="max-w-[120px]">Status</Select.Trigger>
        <Select.Content>
          {#each mediaStatus as status (status.value)}
            <Select.Item value={status.value} label={status.label} />
          {/each}
        </Select.Content>
      </Select.Root>

      <Select.Root type="multiple" bind:value={$form.excludes.rating}>
        <Select.Trigger class="max-w-[120px]">Rating</Select.Trigger>
        <Select.Content>
          {#each mediaRatings as rating (rating.value)}
            <Select.Item value={rating.value} label={rating.label} />
          {/each}
        </Select.Content>
      </Select.Root>
    </div>

    <Select.Root type="single" allowDeselect={false} bind:value={$form.sort}>
      <Select.Trigger class="max-w-[120px]">Sort</Select.Trigger>
      <Select.Content>
        {#each sortTypes as ty (ty.value)}
          <Select.Item value={ty.value} label={ty.label} />
        {/each}
      </Select.Content>
    </Select.Root>

    <Button
      onclick={() => {
        reset({
          data: {
            query: "",
            filter: { type: [], status: [], rating: [] },
            excludes: { type: [], status: [], rating: [] },
            sort: defaultSort,
          },
        });
      }}>Reset filters</Button
    >
  </div>

  <Button type="submit">Filter</Button>
</form>
