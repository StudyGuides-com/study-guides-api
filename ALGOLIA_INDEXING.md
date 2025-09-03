# Feature Overview

This feature will allow clients of the API to use natural language to start the indexing or reindexing of
Tags. 

## Workflow
 1) Outbox population
    - whenever a Tag is created, updated, or deleted, a row is inserted into IndexOutbox.
    - action = upsert for create/update.
    - action = delete for delete.

 2) Client triggers indexing
    - Client app sends a natural language prompt to the API (e.g., "please reindex tags").
    - API uses intent classification / mapping rules to interpret the request into a structured action (e.g., ReindexTags).

 3) Immediate API response
    - API responds with a 202 Accepted or equivalent, acknowledging that indexing has been scheduled.
    - Response includes a job ID or message for tracking.

 4) Background indexing process
    - Worker service reads from IndexOutbox in batches, canonicalizes Tag payloads, computes hashes, and compares against SearchIndexState.
    - Only Tags with new or changed hashes are pushed to Algolia.
    - Deletes are pushed via deleteObjects.

 5) State updates
    - On successful push: update SearchIndexState with the new hash and timestamp.
    - On error: update lastError and retry with exponential backoff.


## Database Schema

```
enum IndexAction {
  upsert
  delete
}

enum IndexObjectType {
  Tag
}

model IndexOutbox {
  id         BigInt          @id @default(autoincrement())
  objectType IndexObjectType
  objectId   String
  action     IndexAction
  queuedAt   DateTime        @default(now())

  @@index([objectType, objectId])
  @@index([queuedAt, id]) // helps worker scan newest-first or oldest-first
}

model SearchIndexState {
  objectType       IndexObjectType
  objectId         String
  lastIndexedAt    DateTime?
  lastIndexedHash  Bytes?
  lastAttemptAt    DateTime?
  attemptCount     Int            @default(0)
  lastError        String?        @db.Text

  @@id([objectType, objectId])
  @@index([lastIndexedAt]) // handy for reconciler windows
}````



### Tag to Algolia Record Transformation

```
import { Prisma } from "@prisma/client";
import { AlgoliaTagRecord } from "./AlgoliaTagRecord";
import { TagInfo, Tag } from "@/types";

export const toAlgoliaTagRecord = (
  tag: Tag,
  accessList: string[],
  ancestry: TagInfo[]
): Prisma.JsonObject => {
  const record: AlgoliaTagRecord = {
    objectID: tag.id,
    id: tag.id,
    name: tag.name,
    description: tag.description,
    hasQuestions: tag.hasQuestions,
    type: tag.type,
    context: tag.context,
    public: tag.public,
    batchId: tag.batchId,
    metaTags: tag.metaTags,
    contentRating: tag.contentRating,
    contentDescriptors: tag.contentDescriptors,
    missingContentDescriptors: tag.contentDescriptors.length === 0,
    missingContentRating: tag.contentRating === "RatingPending",
    missingMetaTags: tag.metaTags.length === 0,
    ownerId: tag.ownerId,
    accessList: accessList,
    tags: ancestry,
  };
  return JSON.parse(JSON.stringify(record));


### Creating record

```
```
```
```
import { Tag } from "@/types";
import { prisma } from "@/lib/database";
import { getTagAncestry } from "@/lib/database/tags/public";
import { getAccessList } from "./getAccessList";
import { toAlgoliaTagRecordType } from "./toAlgoliaTagRecordType";
import { toAlgoliaTagRecord } from "./toAlgoliaTagRecord";

export const createTagRecord = async (tag: Tag) => {
  const accessList = await getAccessList(tag.id);
  const ancestry = await getTagAncestry(tag.id);
  const record = await prisma.algoliaRecord.create({
    data: {
      id: tag.id,
      type: toAlgoliaTagRecordType(tag.type),
      record: toAlgoliaTagRecord(tag, accessList, ancestry),
      createdAt: new Date(),
      updatedAt: new Date(),
      uploaded: false,
    },
  });

  return record;
};```


```import { prisma } from "@/lib/database";

export const getAccessList = async (tagId: string): Promise<string[]> => {
    const accessList = await prisma.tagAccess.findMany({
      where: {
        tagId: tagId,
      },
    });
  
    return accessList.map((access) => access.userId);
  };


import { prisma } from "@/lib/database";
import { TagInfo, Tag } from "@/types";

const toTagInfo = (tag: Tag): TagInfo => ({
  id: tag.id,
  name: tag.name,
  type: tag.type,
  parentTagId: tag.parentTagId,
  hasQuestions: tag.hasQuestions,
  hasChildren: tag.hasChildren,
});

const climbAncestryTree = async (
  currentTagId: string,
  tagAncestry: TagInfo[]
): Promise<void> => {
  const tag = await prisma.tag.findUnique({
    where: { id: currentTagId },
  });

  if (!tag) {
    throw new Error(`Tag with ID ${currentTagId} not found.`);
  }

  const tagWithProperties: Tag = {
    ...tag,
    recentTags: [],
    favoriteTags: [],
    parentTags: [],
    siblingTags: [],
    childTags: [],
  };

  tagAncestry.push(toTagInfo(tagWithProperties));
  if (tag.parentTagId) {
    await climbAncestryTree(tag.parentTagId, tagAncestry);
  }
};

const getTagAncestry = async (tagId: string): Promise<TagInfo[]> => {
  try {
    if (!tagId) {
      return [];
    }

    const tagAncestry: TagInfo[] = [];
    await climbAncestryTree(tagId, tagAncestry);
    return tagAncestry;
  } catch (error) {
    console.error(
      `Error occurred while fetching tag ancestry for ${tagId}:`,
      error
    );
    throw error;
  }
};

export default getTagAncestry;


```
```
