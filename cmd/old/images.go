package old

////func sharePhoto(c *ishell.Context) {
////	if len(c.Args) == 0 {
////		c.Err(errors.New("missing photo id"))
////		return
////	}
////	if len(c.Args) == 1 {
////		c.Err(errors.New("missing destination thread id"))
////		return
////	}
////	id := c.Args[0]
////	threadId := c.Args[1]
////
////	c.Print("caption (optional): ")
////	caption := c.ReadLine()
////
////	// get the original block
////	block, err := getPhotoBlockByDataId(id)
////	if err != nil {
////		c.Err(err)
////		return
////	}
////
////	// lookup destination thread
////	toThread := core.Node.Thread(threadId)
////	if toThread == nil {
////		c.Err(errors.New(fmt.Sprintf("could not find thread %s", threadId)))
////		return
////	}
////
////	// finally, add to destination
////	if _, err := toThread.AddFile(id, caption, block.DataKey); err != nil {
////		c.Err(err)
////		return
////	}
////}
//
////func listPhotos(c *ishell.Context) {
////	if len(c.Args) == 0 {
////		c.Err(errors.New("missing thread id"))
////		return
////	}
////	threadId := c.Args[0]
////
////	thrd := core.Node.Thread(threadId)
////	if thrd == nil {
////		c.Err(errors.New(fmt.Sprintf("could not find thread: %s", threadId)))
////		return
////	}
////
////	query := fmt.Sprintf("threadId='%s' and type=%d", thrd.Id, repo.FilesBlock)
////	blocks := core.Node.Blocks("", -1, query)
////	if len(blocks) == 0 {
////		c.Println(fmt.Sprintf("no photos found in: %s", thrd.Id))
////	} else {
////		c.Println(fmt.Sprintf("%v photos:", len(blocks)))
////	}
////
////	magenta := color.New(color.FgHiMagenta).SprintFunc()
////	for _, block := range blocks {
////		c.Println(magenta(fmt.Sprintf("id: %s, block: %s", block.DataId, block.Id)))
////	}
////}
////
////func getPhoto(c *ishell.Context) {
////	if len(c.Args) == 0 {
////		c.Err(errors.New("missing photo id"))
////		return
////	}
////	if len(c.Args) == 1 {
////		c.Err(errors.New("missing out directory"))
////		return
////	}
////	id := c.Args[0]
////
////	// try to get path with home dir tilda
////	dest, err := homedir.Expand(c.Args[1])
////	if err != nil {
////		dest = c.Args[1]
////	}
////
////	block, err := getPhotoBlockByDataId(id)
////	if err != nil {
////		c.Err(err)
////		return
////	}
////
////	data, err := core.Node.BlockData(fmt.Sprintf("%s/photo", id), block)
////	if err != nil {
////		c.Err(err)
////		return
////	}
////
////	path := filepath.Join(dest, id)
////	if err := ioutil.WriteFile(path, data, 0644); err != nil {
////		c.Err(err)
////		return
////	}
////
////	blue := color.New(color.FgHiBlue).SprintFunc()
////	c.Println(blue("saved to " + path))
////}
////
////func getPhotoMetadata(c *ishell.Context) {
////	if len(c.Args) == 0 {
////		c.Err(errors.New("missing photo id"))
////		return
////	}
////	id := c.Args[0]
////
////	block, err := getPhotoBlockByDataId(id)
////	if err != nil {
////		c.Err(err)
////		return
////	}
////
////	jsonb, err := json.MarshalIndent(block.DataMetadata, "", "    ")
////	if err != nil {
////		c.Err(err)
////		return
////	}
////
////	black := color.New(color.FgHiBlack).SprintFunc()
////	c.Println(black(string(jsonb)))
////}
////
////func getPhotoKey(c *ishell.Context) {
////	if len(c.Args) == 0 {
////		c.Err(errors.New("missing photo id"))
////		return
////	}
////	id := c.Args[0]
////
////	block, err := getPhotoBlockByDataId(id)
////	if err != nil {
////		c.Err(err)
////		return
////	}
////
////	blue := color.New(color.FgHiBlue).SprintFunc()
////	c.Println(blue(base58.FastBase58Encoding(block.DataKey)))
////}
////
////func addPhotoComment(c *ishell.Context) {
////	if len(c.Args) == 0 {
////		c.Err(errors.New("missing block id"))
////		return
////	}
////	id := c.Args[0]
////	c.Print("comment: ")
////	body := c.ReadLine()
////
////	block, err := core.Node.Block(id)
////	if err != nil {
////		c.Err(err)
////		return
////	}
////	thrd := core.Node.Thread(block.ThreadId)
////	if thrd == nil {
////		c.Err(errors.New(fmt.Sprintf("could not find thread %s", block.ThreadId)))
////		return
////	}
////
////	if _, err := thrd.AddComment(block.Id, body); err != nil {
////		c.Err(err)
////		return
////	}
////}
////
////func addPhotoLike(c *ishell.Context) {
////	if len(c.Args) == 0 {
////		c.Err(errors.New("missing block id"))
////		return
////	}
////	id := c.Args[0]
////
////	block, err := core.Node.Block(id)
////	if err != nil {
////		c.Err(err)
////		return
////	}
////	thrd := core.Node.Thread(block.ThreadId)
////	if thrd == nil {
////		c.Err(errors.New(fmt.Sprintf("could not find thread %s", block.ThreadId)))
////		return
////	}
////
////	if _, err := thrd.AddLike(block.Id); err != nil {
////		c.Err(err)
////		return
////	}
////}
////
////func listPhotoComments(c *ishell.Context) {
////	if len(c.Args) == 0 {
////		c.Err(errors.New("missing block id"))
////		return
////	}
////	id := c.Args[0]
////
////	block, err := core.Node.Block(id)
////	if err != nil {
////		c.Err(err)
////		return
////	}
////	thrd := core.Node.Thread(block.ThreadId)
////	if thrd == nil {
////		c.Err(errors.New(fmt.Sprintf("could not find thread %s", block.ThreadId)))
////		return
////	}
////
////	query := fmt.Sprintf("threadId='%s' and type=%d", thrd.Id, repo.CommentBlock)
////	blocks := core.Node.Blocks("", -1, query)
////	if len(blocks) == 0 {
////		c.Println(fmt.Sprintf("no comments found on: %s", block.Id))
////	} else {
////		c.Println(fmt.Sprintf("%v comments:", len(blocks)))
////	}
////
////	cyan := color.New(color.FgHiCyan).SprintFunc()
////	for _, b := range blocks {
////		username := core.Node.ContactUsername(b.AuthorId)
////		c.Println(cyan(fmt.Sprintf("%s: %s: %s", b.Id, username, b.DataCaption)))
////	}
////}
////
////func listPhotoLikes(c *ishell.Context) {
////	if len(c.Args) == 0 {
////		c.Err(errors.New("missing block id"))
////		return
////	}
////	id := c.Args[0]
////
////	block, err := core.Node.Block(id)
////	if err != nil {
////		c.Err(err)
////		return
////	}
////	thrd := core.Node.Thread(block.ThreadId)
////	if thrd == nil {
////		c.Err(errors.New(fmt.Sprintf("could not find thread %s", block.ThreadId)))
////		return
////	}
////
////	query := fmt.Sprintf("threadId='%s' and type=%d", thrd.Id, repo.LikeBlock)
////	blocks := core.Node.Blocks("", -1, query)
////	if len(blocks) == 0 {
////		c.Println(fmt.Sprintf("no likes found on: %s", block.Id))
////	} else {
////		c.Println(fmt.Sprintf("%v likes:", len(blocks)))
////	}
////
////	cyan := color.New(color.FgHiCyan).SprintFunc()
////	for _, b := range blocks {
////		username := core.Node.ContactUsername(b.AuthorId)
////		c.Println(cyan(fmt.Sprintf("%s: %s", b.Id, username)))
////	}
////}
////
////func getPhotoBlockByDataId(dataId string) (*repo.Block, error) {
////	block, err := core.Node.BlockByDataId(dataId)
////	if err != nil {
////		return nil, err
////	}
////	if block.Type != repo.FilesBlock {
////		return nil, errors.New("not a photo block, aborting")
////	}
////	return block, nil
////}
