package keeper_test

import (
	"testing"

	"github.com/stretchr/testify/suite"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"
	tmtime "github.com/tendermint/tendermint/types/time"

	"github.com/cosmos/cosmos-sdk/simapp"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/nft"
)

const (
	testClassID          = "kitty"
	testClassName        = "Crypto Kitty"
	testClassSymbol      = "kitty"
	testClassDescription = "Crypto Kitty"
	testClassURI         = "class uri"
	testClassURIHash     = "ae702cefd6b6a65fe2f991ad6d9969ed"
	testID               = "kitty1"
	testURI              = "kitty uri"
	testURIHash          = "229bfd3c1b431c14a526497873897108"
)

type TestSuite struct {
	suite.Suite

	app         *simapp.SimApp
	ctx         sdk.Context
	addrs       []sdk.AccAddress
	queryClient nft.QueryClient
}

func (s *TestSuite) SetupTest() {
	app := simapp.Setup(s.T(), false)
	ctx := app.BaseApp.NewContext(false, tmproto.Header{})
	ctx = ctx.WithBlockHeader(tmproto.Header{Time: tmtime.Now()})

	s.app = app
	s.ctx = ctx
	s.addrs = simapp.AddTestAddrsIncremental(app, ctx, 3, sdk.NewInt(30000000))
}

func TestTestSuite(t *testing.T) {
	suite.Run(t, new(TestSuite))
}

func (s *TestSuite) TestSaveClass() {
	except := nft.Class{
		Id:          testClassID,
		Name:        testClassName,
		Symbol:      testClassSymbol,
		Description: testClassDescription,
		Uri:         testClassURI,
		UriHash:     testClassURIHash,
	}
	err := s.app.NFTKeeper.SaveClass(s.ctx, except)
	s.Require().NoError(err)

	actual, has := s.app.NFTKeeper.GetClass(s.ctx, testClassID)
	s.Require().True(has)
	s.Require().EqualValues(except, actual)

	classes := s.app.NFTKeeper.GetClasses(s.ctx)
	s.Require().EqualValues([]*nft.Class{&except}, classes)
}

func (s *TestSuite) TestUpdateClass() {
	class := nft.Class{
		Id:          testClassID,
		Name:        testClassName,
		Symbol:      testClassSymbol,
		Description: testClassDescription,
		Uri:         testClassURI,
		UriHash:     testClassURIHash,
	}
	err := s.app.NFTKeeper.SaveClass(s.ctx, class)
	s.Require().NoError(err)

	noExistClass := nft.Class{
		Id:          "kitty1",
		Name:        testClassName,
		Symbol:      testClassSymbol,
		Description: testClassDescription,
		Uri:         testClassURI,
		UriHash:     testClassURIHash,
	}

	err = s.app.NFTKeeper.UpdateClass(s.ctx, noExistClass)
	s.Require().Error(err)
	s.Require().Contains(err.Error(), "nft class does not exist")

	except := nft.Class{
		Id:          testClassID,
		Name:        "My crypto Kitty",
		Symbol:      testClassSymbol,
		Description: testClassDescription,
		Uri:         testClassURI,
		UriHash:     testClassURIHash,
	}

	err = s.app.NFTKeeper.UpdateClass(s.ctx, except)
	s.Require().NoError(err)

	actual, has := s.app.NFTKeeper.GetClass(s.ctx, testClassID)
	s.Require().True(has)
	s.Require().EqualValues(except, actual)
}

func (s *TestSuite) TestMint() {
	class := nft.Class{
		Id:          testClassID,
		Name:        testClassName,
		Symbol:      testClassSymbol,
		Description: testClassDescription,
		Uri:         testClassURI,
		UriHash:     testClassURIHash,
	}
	err := s.app.NFTKeeper.SaveClass(s.ctx, class)
	s.Require().NoError(err)

	expNFT := nft.NFT{
		ClassId: testClassID,
		Id:      testID,
		Uri:     testURI,
	}
	err = s.app.NFTKeeper.Mint(s.ctx, expNFT, s.addrs[0])
	s.Require().NoError(err)

	// test GetNFT
	actNFT, has := s.app.NFTKeeper.GetNFT(s.ctx, testClassID, testID)
	s.Require().True(has)
	s.Require().EqualValues(expNFT, actNFT)

	// test GetOwner
	owner := s.app.NFTKeeper.GetOwner(s.ctx, testClassID, testID)
	s.Require().True(s.addrs[0].Equals(owner))

	// test GetNFTsOfClass
	actNFTs := s.app.NFTKeeper.GetNFTsOfClass(s.ctx, testClassID)
	s.Require().EqualValues([]nft.NFT{expNFT}, actNFTs)

	// test GetNFTsOfClassByOwner
	actNFTs = s.app.NFTKeeper.GetNFTsOfClassByOwner(s.ctx, testClassID, s.addrs[0])
	s.Require().EqualValues([]nft.NFT{expNFT}, actNFTs)

	// test GetBalance
	balance := s.app.NFTKeeper.GetBalance(s.ctx, testClassID, s.addrs[0])
	s.Require().EqualValues(uint64(1), balance)

	// test GetTotalSupply
	supply := s.app.NFTKeeper.GetTotalSupply(s.ctx, testClassID)
	s.Require().EqualValues(uint64(1), supply)
}

func (s *TestSuite) TestBurn() {
	except := nft.Class{
		Id:          testClassID,
		Name:        testClassName,
		Symbol:      testClassSymbol,
		Description: testClassDescription,
		Uri:         testClassURI,
		UriHash:     testClassURIHash,
	}
	err := s.app.NFTKeeper.SaveClass(s.ctx, except)
	s.Require().NoError(err)

	expNFT := nft.NFT{
		ClassId: testClassID,
		Id:      testID,
		Uri:     testURI,
	}
	err = s.app.NFTKeeper.Mint(s.ctx, expNFT, s.addrs[0])
	s.Require().NoError(err)

	err = s.app.NFTKeeper.Burn(s.ctx, testClassID, testID)
	s.Require().NoError(err)

	// test GetNFT
	_, has := s.app.NFTKeeper.GetNFT(s.ctx, testClassID, testID)
	s.Require().False(has)

	// test GetOwner
	owner := s.app.NFTKeeper.GetOwner(s.ctx, testClassID, testID)
	s.Require().Nil(owner)

	// test GetNFTsOfClass
	actNFTs := s.app.NFTKeeper.GetNFTsOfClass(s.ctx, testClassID)
	s.Require().Empty(actNFTs)

	// test GetNFTsOfClassByOwner
	actNFTs = s.app.NFTKeeper.GetNFTsOfClassByOwner(s.ctx, testClassID, s.addrs[0])
	s.Require().Empty(actNFTs)

	// test GetBalance
	balance := s.app.NFTKeeper.GetBalance(s.ctx, testClassID, s.addrs[0])
	s.Require().EqualValues(uint64(0), balance)

	// test GetTotalSupply
	supply := s.app.NFTKeeper.GetTotalSupply(s.ctx, testClassID)
	s.Require().EqualValues(uint64(0), supply)
}

func (s *TestSuite) TestUpdate() {
	class := nft.Class{
		Id:          testClassID,
		Name:        testClassName,
		Symbol:      testClassSymbol,
		Description: testClassDescription,
		Uri:         testClassURI,
		UriHash:     testClassURIHash,
	}
	err := s.app.NFTKeeper.SaveClass(s.ctx, class)
	s.Require().NoError(err)

	myNFT := nft.NFT{
		ClassId: testClassID,
		Id:      testID,
		Uri:     testURI,
	}
	err = s.app.NFTKeeper.Mint(s.ctx, myNFT, s.addrs[0])
	s.Require().NoError(err)

	expNFT := nft.NFT{
		ClassId: testClassID,
		Id:      testID,
		Uri:     "updated",
	}

	err = s.app.NFTKeeper.Update(s.ctx, expNFT)
	s.Require().NoError(err)

	// test GetNFT
	actNFT, has := s.app.NFTKeeper.GetNFT(s.ctx, testClassID, testID)
	s.Require().True(has)
	s.Require().EqualValues(expNFT, actNFT)
}

func (s *TestSuite) TestTransfer() {
	class := nft.Class{
		Id:          testClassID,
		Name:        testClassName,
		Symbol:      testClassSymbol,
		Description: testClassDescription,
		Uri:         testClassURI,
		UriHash:     testClassURIHash,
	}
	err := s.app.NFTKeeper.SaveClass(s.ctx, class)
	s.Require().NoError(err)

	expNFT := nft.NFT{
		ClassId: testClassID,
		Id:      testID,
		Uri:     testURI,
	}
	err = s.app.NFTKeeper.Mint(s.ctx, expNFT, s.addrs[0])
	s.Require().NoError(err)

	//valid owner
	err = s.app.NFTKeeper.Transfer(s.ctx, testClassID, testID, s.addrs[1])
	s.Require().NoError(err)

	// test GetOwner
	owner := s.app.NFTKeeper.GetOwner(s.ctx, testClassID, testID)
	s.Require().Equal(s.addrs[1], owner)

	balanceAddr0 := s.app.NFTKeeper.GetBalance(s.ctx, testClassID, s.addrs[0])
	s.Require().EqualValues(uint64(0), balanceAddr0)

	balanceAddr1 := s.app.NFTKeeper.GetBalance(s.ctx, testClassID, s.addrs[1])
	s.Require().EqualValues(uint64(1), balanceAddr1)

	// test GetNFTsOfClassByOwner
	actNFTs := s.app.NFTKeeper.GetNFTsOfClassByOwner(s.ctx, testClassID, s.addrs[1])
	s.Require().EqualValues([]nft.NFT{expNFT}, actNFTs)
}

func (s *TestSuite) TestExportGenesis() {
	class := nft.Class{
		Id:          testClassID,
		Name:        testClassName,
		Symbol:      testClassSymbol,
		Description: testClassDescription,
		Uri:         testClassURI,
		UriHash:     testClassURIHash,
	}
	err := s.app.NFTKeeper.SaveClass(s.ctx, class)
	s.Require().NoError(err)

	expNFT := nft.NFT{
		ClassId: testClassID,
		Id:      testID,
		Uri:     testURI,
	}
	err = s.app.NFTKeeper.Mint(s.ctx, expNFT, s.addrs[0])
	s.Require().NoError(err)

	expGenesis := &nft.GenesisState{
		Classes: []*nft.Class{&class},
		Entries: []*nft.Entry{{
			Owner: s.addrs[0].String(),
			Nfts:  []*nft.NFT{&expNFT},
		}},
	}
	genesis := s.app.NFTKeeper.ExportGenesis(s.ctx)
	s.Require().Equal(expGenesis, genesis)
}

func (s *TestSuite) TestInitGenesis() {
	expClass := nft.Class{
		Id:          testClassID,
		Name:        testClassName,
		Symbol:      testClassSymbol,
		Description: testClassDescription,
		Uri:         testClassURI,
		UriHash:     testClassURIHash,
	}
	expNFT := nft.NFT{
		ClassId: testClassID,
		Id:      testID,
		Uri:     testURI,
	}
	expGenesis := &nft.GenesisState{
		Classes: []*nft.Class{&expClass},
		Entries: []*nft.Entry{{
			Owner: s.addrs[0].String(),
			Nfts:  []*nft.NFT{&expNFT},
		}},
	}
	s.app.NFTKeeper.InitGenesis(s.ctx, expGenesis)

	actual, has := s.app.NFTKeeper.GetClass(s.ctx, testClassID)
	s.Require().True(has)
	s.Require().EqualValues(expClass, actual)

	// test GetNFT
	actNFT, has := s.app.NFTKeeper.GetNFT(s.ctx, testClassID, testID)
	s.Require().True(has)
	s.Require().EqualValues(expNFT, actNFT)
}
