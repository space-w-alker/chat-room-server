-- CreateTable
CREATE TABLE "RoomTag" (
    "id" TEXT NOT NULL,
    "name" TEXT NOT NULL,

    CONSTRAINT "RoomTag_pkey" PRIMARY KEY ("id")
);

-- CreateTable
CREATE TABLE "_ChatRoomToRoomTag" (
    "A" TEXT NOT NULL,
    "B" TEXT NOT NULL
);

-- CreateIndex
CREATE UNIQUE INDEX "_ChatRoomToRoomTag_AB_unique" ON "_ChatRoomToRoomTag"("A", "B");

-- CreateIndex
CREATE INDEX "_ChatRoomToRoomTag_B_index" ON "_ChatRoomToRoomTag"("B");

-- AddForeignKey
ALTER TABLE "_ChatRoomToRoomTag" ADD CONSTRAINT "_ChatRoomToRoomTag_A_fkey" FOREIGN KEY ("A") REFERENCES "ChatRoom"("id") ON DELETE CASCADE ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE "_ChatRoomToRoomTag" ADD CONSTRAINT "_ChatRoomToRoomTag_B_fkey" FOREIGN KEY ("B") REFERENCES "RoomTag"("id") ON DELETE CASCADE ON UPDATE CASCADE;
