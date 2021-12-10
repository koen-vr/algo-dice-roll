from pyteal import *

# Demo's an on-chain dice roller with a pseudo random number generator.
# It is left up to the user of the contract to decide the type of dice
# are represented by the byte values output by the roller application.
def approval_program():

    idx = ScratchVar(TealType.uint64)
    val = ScratchVar(TealType.uint64)

    state = ScratchVar(TealType.uint64)

    hand = ScratchVar(TealType.bytes)
    result = ScratchVar(TealType.bytes)

    get_next = Seq([
        state.store(state.load() ^ (state.load() >> Int(12))),
        state.store(state.load() ^ (state.load() << Int(25))),
        state.store(state.load() ^ (state.load() >> Int(27))),
        result.store(BytesMul(Itob(state.load()), Bytes("base16", "0x2545F4914F6CDD1D"))),
        result.store(Extract(result.load(), Len(result.load())-Int(8), Int(8)))
    ])

    set_seed = Seq([
        # Byte padding
        result.store(BytesAdd(Itob(val.load()), Bytes("base16", "0x9E3779B97f4A7C15"))),
        val.store(Btoi(Extract(result.load(), Len(result.load())-Int(8), Int(8)))),
        ## First shift round
        val.store(val.load() ^ (val.load() >> Int(30))),
        result.store(BytesMul(Itob(val.load()), Bytes("base16", "0xBF58476D1CE4E5B9"))),
        val.store(Btoi(Extract(result.load(), Len(result.load())-Int(8), Int(8)))),
        ## Second shift round
        val.store(val.load() ^ (val.load() >> Int(27))),
        result.store(BytesMul(Itob(val.load()), Bytes("base16", "0x94D049BB133111EB"))),
        val.store(Btoi(Extract(result.load(), Len(result.load())-Int(8), Int(8)))),
        ## Equalize the weights
        state.store(val.load() ^ (val.load() >> Int(31))),
    ])

    set_hand = Seq([
        Assert(Len(result.load()) >= Int(8)),
        hand.store(Bytes("base16", "0xFFFFFFFFFFFFFFFF")),
        For(idx.store(Int(0)), idx.load() < Int(8), idx.store(idx.load() + Int(1))).Do(
            Seq([
                val.store(GetByte(result.load(), idx.load())),
                val.store(Int(1)+(((val.load()*Int(10)) /Int(256)) *Int(6)) /Int(10)),
                hand.store(SetByte(hand.load(), idx.load(), val.load())),
            ])
        )
    ])

    seed_hand = Seq([
        # Load A - args['gm-seed']
        # Load B - args['act-seed']
        # Load C - glob['xor-seed']
        # A is Split A XOR C
        # C is Split B XOR A
        # State is Shift C
        # Store State
    ])

    roll_hand = Seq([
        # Load State
        # Shift State
        # Store State
        # Multiply result
        # Store the result
    ])

    handle_noop = Seq([
        Return(Int(1))
    ])

    handle_optin = Seq([
        Return(Int(1))
    ])

    handle_create = Seq([
        val.store(Int(123456258741789)),
        set_seed,

        get_next,
        App.globalPut(Bytes("r1"), Btoi(result.load())),
        set_hand,
        App.globalPut(Bytes("h1"), hand.load()),

        Return(Int(1))
    ])

    handle_closeout = Seq([
        Return(Int(1))
    ])

    handle_updateapp = Err()

    handle_deleteapp = Err()

    program = Cond(
        [Txn.application_id() == Int(0), handle_create],
        [Txn.on_completion() == OnComplete.NoOp, handle_noop],
        [Txn.on_completion() == OnComplete.OptIn, handle_optin],
        [Txn.on_completion() == OnComplete.CloseOut, handle_closeout],
        [Txn.on_completion() == OnComplete.UpdateApplication, handle_updateapp],
        [Txn.on_completion() == OnComplete.DeleteApplication, handle_deleteapp]
    )
    return program

if __name__ == "__main__":
    print(compileTeal(approval_program(), Mode.Application, version=5))