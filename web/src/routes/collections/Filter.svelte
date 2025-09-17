<script lang="ts">
  import Errors from "$lib/components/Errors.svelte";
  import FormItem from "$lib/components/FormItem.svelte";
  import { Button, Input, Label, Select } from "@nanoteck137/nano-ui";
  import { zod } from "sveltekit-superforms/adapters";
  import { defaults, superForm } from "sveltekit-superforms/client";
  import { FullFilter } from "./types";
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
      query.delete("excludes");
      query.delete("filter");

      query.set("query", data.query);
      query.set("sort", data.sort);

      if (data.excludes.length > 0) {
        query.delete("filter");
        query.set("excludes", data.excludes.join(","));
      }

      if (data.filter !== "") {
        query.delete("excludes");
        query.set("filter", data.filter);
      }

      goto("?" + query.toString(), { invalidateAll: true });
    }, 0);
  }

  const { form, errors, enhance } = superForm(
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
    <Label for="query">Search Collections</Label>
    <Input id="query" name="query" type="text" bind:value={$form.query} />
    <Errors errors={$errors.query} />
  </FormItem>

  <div class="flex gap-4">
    <Select.Root
      type="single"
      bind:value={$form.filter}
      disabled={$form.excludes.length > 0}
    >
      <Select.Trigger class="max-w-[120px]">Filter</Select.Trigger>
      <Select.Content>
        <Select.Item value="unknown" label="Unknown" />
        <Select.Item value="series" label="Series" />
        <Select.Item value="anime" label="Anime" />
      </Select.Content>
    </Select.Root>

    <Select.Root
      type="multiple"
      bind:value={$form.excludes}
      disabled={$form.filter !== ""}
    >
      <Select.Trigger class="max-w-[120px]">Exclude</Select.Trigger>
      <Select.Content>
        <Select.Item value="unknown" label="Unknown" />
        <Select.Item value="series" label="Series" />
        <Select.Item value="anime" label="Anime" />
      </Select.Content>
    </Select.Root>

    <Select.Root type="single" allowDeselect={false} bind:value={$form.sort}>
      <Select.Trigger class="max-w-[120px]">Sort</Select.Trigger>
      <Select.Content>
        <Select.Item value="name-a-z" label={"Name (A-Z)"} />
        <Select.Item value="name-z-a" label={"Name (A-Z)"} />
      </Select.Content>
    </Select.Root>
  </div>

  <Button type="submit">Filter</Button>
</form>
