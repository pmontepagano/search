package ar.com.montepagano.search.examples.creditcardpayments;

import java.io.FileInputStream;
import java.io.IOException;
import java.util.Scanner;

import ar.com.montepagano.search.middleware.v1.Middleware;
import org.json.JSONObject;
import io.grpc.ManagedChannel;
import io.grpc.ManagedChannelBuilder;
import com.google.protobuf.ByteString;
import ar.com.montepagano.search.middleware.v1.PrivateMiddlewareServiceGrpc;
import ar.com.montepagano.search.middleware.v1.Middleware.RegisterChannelRequest;
import ar.com.montepagano.search.middleware.v1.Middleware.RegisterChannelResponse;
import ar.com.montepagano.search.appmessage.v1.AppMessageOuterClass.AppSendRequest;
import ar.com.montepagano.search.appmessage.v1.AppMessageOuterClass.AppMessage;
import ar.com.montepagano.search.contracts.v1.Contracts.GlobalContract;
import ar.com.montepagano.search.contracts.v1.Contracts.GlobalContractFormat;
import ar.com.montepagano.search.broker.v1.Broker.RemoteParticipant;

public class Main {
    public static void main(String[] args) {
        Scanner scanner = new Scanner(System.in);

        // List of book titles
        String[] bookTitles = {
                "The Great Gatsby",
                "To Kill a Mockingbird",
                "1984",
                "Pride and Prejudice",
                "Animal Farm",
                "The Catcher in the Rye",
                "Lord of the Flies",
                "Brave New World",
                "The Hobbit",
                "The Fellowship of the Ring"
        };

        // Prompt user to select book titles
        System.out.println("Please select at least one book to purchase:");
        for (int i = 0; i < bookTitles.length; i++) {
            System.out.println((i + 1) + ". " + bookTitles[i]);
        }

        // Get user input for selected book titles
        System.out.print("Enter the number(s) of the book(s) you want to purchase (comma-separated): ");
        String selectedBooksInput = scanner.nextLine();
        String[] selectedBooksInputArray = selectedBooksInput.split(",");
        int[] selectedBooks = new int[selectedBooksInputArray.length];
        for (int i = 0; i < selectedBooksInputArray.length; i++) {
            selectedBooks[i] = Integer.parseInt(selectedBooksInputArray[i].trim());
        }

        // Prompt user for shipping address
        System.out.print("Enter your shipping address: ");
        String shippingAddress = scanner.nextLine();

        // Print purchase information
        System.out.println("Purchase Information:");
        System.out.println("Selected Books:");
        for (int selectedBook : selectedBooks) {
            System.out.println("- " + bookTitles[selectedBook - 1]);
        }
        System.out.println("Shipping Address: " + shippingAddress);

        // Load file contract.fsa into a GlobalContract
        ByteString contractBytes = null;
        try {
            contractBytes = ByteString.readFrom(new FileInputStream("./contract.fsa"));
        } catch (IOException e) {
            e.printStackTrace();
        }
        GlobalContract contract = GlobalContract.newBuilder().setContract(contractBytes).setFormat(
                GlobalContractFormat.GLOBAL_CONTRACT_FORMAT_FSA).setInitiatorName("ClientApp").build();

        // Get the stub to communicate with the middleware
        ManagedChannel channel = ManagedChannelBuilder.forTarget("middleware-client:11000").usePlaintext().build();
        PrivateMiddlewareServiceGrpc.PrivateMiddlewareServiceBlockingStub stub = PrivateMiddlewareServiceGrpc.newBlockingStub(channel);

        // Register channel with middleware using GlobalContract.
        // Add preset participants to the RegisterChannelRequest. We do this because we don't have an algorithm for compatibility checking in the broker.
        RemoteParticipant pps = RemoteParticipant.newBuilder().setUrl("middleware-payments:10000").setAppId("PPS").build();
        RemoteParticipant srv = RemoteParticipant.newBuilder().setUrl("middleware-backend:10000").setAppId("Srv").build();
        RegisterChannelRequest request = RegisterChannelRequest.newBuilder().setRequirementsContract(contract).putPresetParticipants("PPS", pps).putPresetParticipants("Srv", srv).build();
        RegisterChannelResponse response = stub.registerChannel(request);
        var channelId = response.getChannelId();

        // Send PurchaseRequest with each item quantities and the shipping address
        var body = ByteString.copyFromUtf8("{}");
        var msg = AppMessage.newBuilder().setType("PurchaseRequest").setBody(body).build();
        var sendreq = AppSendRequest.newBuilder().setChannelId(channelId).setRecipient("Srv").setMessage(msg).build();
        var sendresp = stub.appSend(sendreq);
        if (sendresp.getResult() != Middleware.AppSendResponse.Result.RESULT_OK) {
            System.out.println("Error sending PurchaseRequest. Exiting...");
            System.exit(1);
        }


        channel.shutdown();
    }
}