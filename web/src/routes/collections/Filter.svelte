<script lang="ts">
  import Errors from "$lib/components/Errors.svelte";
  import FormItem from "$lib/components/FormItem.svelte";
  import { Button, Input, Label, Select } from "@nanoteck137/nano-ui";
  import { zod } from "sveltekit-superforms/adapters";
  import { defaults, superForm } from "sveltekit-superforms/client";
  import {
    collectionTypes,
    defaultSort,
    FullFilter,
    sortTypes,
  } from "./types";
  import { goto } from "$app/navigation";
  import { page } from "$app/stores";

  export type Props = {
    fullFilter: FullFilter;
  };

  const { fullFilter }: Props = $props();

  function submit(data: FullFilter) {
    setTimeout(() => {
      const query = $page.url.searchParams;
      query.delete("query");
      query.delete("sort");

      query.delete("filterType");

      query.delete("excludeType");

      query.set("query", data.query);
      query.set("sort", data.sort);

      if (data.filters.type.length > 0) {
        query.set("filterType", data.filters.type.join(","));
      }

      if (data.excludes.type.length > 0) {
        query.set("excludeType", data.excludes.type.join(","));
      }

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
      <p class="w-20">Filters</p>

      <Select.Root type="multiple" bind:value={$form.filters.type}>
        <Select.Trigger class="max-w-[120px]">Type</Select.Trigger>
        <Select.Content>
          {#each collectionTypes as ty (ty.value)}
            <Select.Item value={ty.value} label={ty.label} />
          {/each}
        </Select.Content>
      </Select.Root>
    </div>

    <div class="flex items-center gap-4">
      <p class="w-20">Excludes</p>

      <Select.Root type="multiple" bind:value={$form.excludes.type}>
        <Select.Trigger class="max-w-[120px]">Type</Select.Trigger>
        <Select.Content>
          {#each collectionTypes as ty (ty.value)}
            <Select.Item value={ty.value} label={ty.label} />
          {/each}
        </Select.Content>
      </Select.Root>
    </div>

    <div class="flex items-center gap-4">
      <p class="w-20">Sort</p>

      <Select.Root type="single" allowDeselect={false} bind:value={$form.sort}>
        <Select.Trigger class="max-w-[160px]">
          {sortTypes.find((i) => i.value === $form.sort)?.label ?? "Sort"}
        </Select.Trigger>
        <Select.Content>
          {#each sortTypes as ty (ty.value)}
            <Select.Item value={ty.value} label={ty.label} />
          {/each}
        </Select.Content>
      </Select.Root>
    </div>

    <Button
      onclick={() => {
        reset({
          data: {
            query: "",
            filters: { type: [] },
            excludes: { type: [] },
            sort: defaultSort,
          },
        });
      }}
    >
      Reset filters
    </Button>

    <Button type="submit">Filter</Button>
  </div>
</form>
